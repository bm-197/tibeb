package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Config holds the configuration for code generation
type Config struct {
	InputFile string
	OutputDir string
	Package   string
	Verbose   bool
}

// ValidationField represents a field in a validation schema
type ValidationField struct {
	Name       string
	Type       string
	Validators []string
}

// ValidationSchema represents a validation schema
type ValidationSchema struct {
	TypeName string
	Fields   []ValidationField
}

// Generate generates validation code from the input file
func Generate(config *Config) error {
	// Parse input file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, config.InputFile, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing input file: %w", err)
	}

	if config.Verbose {
		fmt.Printf("Parsed file: %s\n", config.InputFile)
		ast.Print(fset, f)
	}

	// Find validation schemas
	schemas := findValidationSchemas(f)
	if len(schemas) == 0 {
		return fmt.Errorf("no validation schemas found in %s", config.InputFile)
	}

	// Generate code for each schema
	for _, schema := range schemas {
		if err := generateValidator(config, schema); err != nil {
			return fmt.Errorf("generating validator for %s: %w", schema.TypeName, err)
		}
	}

	return nil
}

// findValidationSchemas looks for validation schema definitions in the AST
func findValidationSchemas(f *ast.File) []ValidationSchema {
	var schemas []ValidationSchema

	ast.Inspect(f, func(n ast.Node) bool {
		// Look for variable declarations that create validation schemas
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for i, value := range valueSpec.Values {
						if schema := extractValidationSchema(value); schema != nil {
							// Try to extract type name from comments or variable name
							if schema.TypeName == "" && i < len(valueSpec.Names) {
								if genDecl.Doc != nil && len(genDecl.Doc.List) > 0 {
									text := genDecl.Doc.List[0].Text
									if strings.Contains(text, "validation schema for") {
										parts := strings.Split(text, "validation schema for")
										if len(parts) > 1 {
											schema.TypeName = strings.TrimSpace(parts[1])
										}
									}
								}
								if schema.TypeName == "" {
									// Try to extract type name from variable name
									varName := valueSpec.Names[i].Name
									if strings.HasSuffix(varName, "Schema") {
										schema.TypeName = strings.TrimSuffix(varName, "Schema")
									} else {
										schema.TypeName = varName
									}
								}
							}
							schemas = append(schemas, *schema)
						}
					}
				}
			}
		}
		return true
	})

	return schemas
}

// extractValidationSchema extracts validation schema from an AST expression
func extractValidationSchema(expr ast.Expr) *ValidationSchema {
	// Look for validate.Struct[Type]() call
	if call, ok := expr.(*ast.CallExpr); ok {
		// Find the root call in the chain
		rootCall := findRootCall(call)
		if rootCall == nil {
			return nil
		}

		// Extract schema type and create schema
		var schema *ValidationSchema
		if sel, ok := rootCall.Fun.(*ast.SelectorExpr); ok {
			if indexExpr, ok := sel.X.(*ast.IndexExpr); ok {
				if rootSel, ok := indexExpr.X.(*ast.SelectorExpr); ok {
					if pkg, ok := rootSel.X.(*ast.Ident); ok && pkg.Name == "validate" && rootSel.Sel.Name == "Struct" {
						if typeIdent, ok := indexExpr.Index.(*ast.Ident); ok {
							schema = &ValidationSchema{
								TypeName: typeIdent.Name,
								Fields:   make([]ValidationField, 0),
							}
						}
					}
				}
			}
		}

		if schema == nil {
			return nil
		}

		// Walk up the chain to collect all Field() calls
		current := call
		for current != nil {
			if sel, ok := current.Fun.(*ast.SelectorExpr); ok {
				if sel.Sel.Name == "Field" {
					field := extractFieldValidation(current)
					if field != nil {
						schema.Fields = append([]ValidationField{*field}, schema.Fields...)
					}
				}
				if callExpr, ok := sel.X.(*ast.CallExpr); ok {
					current = callExpr
				} else {
					break
				}
			} else {
				break
			}
		}

		return schema
	}

	return nil
}

// findRootCall finds the root validate.Struct call in a chain
func findRootCall(call *ast.CallExpr) *ast.CallExpr {
	current := call
	for {
		sel, ok := current.Fun.(*ast.SelectorExpr)
		if !ok {
			break
		}
		callExpr, ok := sel.X.(*ast.CallExpr)
		if !ok {
			break
		}
		current = callExpr
	}
	return current
}

// extractFieldValidation extracts field validation from a Field() call
func extractFieldValidation(call *ast.CallExpr) *ValidationField {
	if len(call.Args) != 2 {
		return nil
	}

	// Extract field name from selector function
	if funcLit, ok := call.Args[0].(*ast.FuncLit); ok {
		if len(funcLit.Body.List) > 0 {
			if returnStmt, ok := funcLit.Body.List[0].(*ast.ReturnStmt); ok {
				if len(returnStmt.Results) > 0 {
					if sel, ok := returnStmt.Results[0].(*ast.SelectorExpr); ok {
						return &ValidationField{
							Name:       sel.Sel.Name,
							Type:       inferFieldType(funcLit.Type.Results),
							Validators: extractValidators(call.Args[1]),
						}
					}
				}
			}
		}
	}
	return nil
}

// inferFieldType infers the field type from the function results
func inferFieldType(results *ast.FieldList) string {
	if results != nil && len(results.List) > 0 {
		if ident, ok := results.List[0].Type.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return "interface{}"
}

// extractValidators extracts validators from a validator chain
func extractValidators(expr ast.Expr) []string {
	var validators []string
	current := expr

	for {
		call, ok := current.(*ast.CallExpr)
		if !ok {
			break
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			break
		}

		// Add validator name
		validators = append([]string{sel.Sel.Name}, validators...)

		// Move to next in chain
		current = sel.X
	}

	return validators
}

// generateValidator generates the validator code for a schema
func generateValidator(config *Config, schema ValidationSchema) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Prepare output file path
	outFile := filepath.Join(config.OutputDir, strings.ToLower(schema.TypeName)+"_validator.go")

	// Create output file
	f, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer f.Close()

	// Parse validator template
	tmpl, err := template.New("validator").Parse(`// Code generated by tibeb. DO NOT EDIT.
package {{ .Package }}

import (
	"github.com/bm-197/tibeb/pkg/validate"
)

// Validate validates the {{ .Schema.TypeName }} struct
func (v {{ .Schema.TypeName }}) Validate() error {
	return {{ .Schema.TypeName }}Schema.Validate(v)
}

// {{ .Schema.TypeName }}Schema is the validation schema for {{ .Schema.TypeName }}
var {{ .Schema.TypeName }}Schema = validate.Struct[{{ .Schema.TypeName }}]().
	{{- range .Schema.Fields }}
	Field(func(v {{ $.Schema.TypeName }}) {{ .Type }} { return v.{{ .Name }} }, validate.{{ .Type }}()
		{{- range .Validators }}.{{ . }}{{ end }})
	{{- end }}
`)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	// Execute template
	data := struct {
		Package string
		Schema  ValidationSchema
	}{
		Package: config.Package,
		Schema:  schema,
	}
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("executing template: %w", err)
	}

	return nil
}
