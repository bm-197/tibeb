package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bm-197/tibeb/internal/generator"
)

func main() {
	genCmd := flag.NewFlagSet("gen", flag.ExitOnError)
	var (
		inputFile string
		outputDir string
		pkgName   string
		verbose   bool
	)

	genCmd.StringVar(&inputFile, "file", "", "Input file containing validation schemas")
	genCmd.StringVar(&outputDir, "out", "", "Output directory for generated code (default: same as input)")
	genCmd.StringVar(&pkgName, "pkg", "", "Package name for generated code (default: same as input)")
	genCmd.BoolVar(&verbose, "verbose", false, "Print verbose output")

	if len(os.Args) < 2 {
		fmt.Println("expected 'gen' subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "gen":
		genCmd.Parse(os.Args[2:])
	default:
		fmt.Printf("unknown subcommand: %s\n", os.Args[1])
		os.Exit(1)
	}

	if inputFile == "" {
		fmt.Fprintln(os.Stderr, "Error: input file is required")
		genCmd.Usage()
		os.Exit(1)
	}

	if outputDir == "" {
		outputDir = filepath.Dir(inputFile)
	}

	if pkgName == "" {
		// Default to the directory name
		pkgName = filepath.Base(filepath.Dir(inputFile))
	}

	config := &generator.Config{
		InputFile: inputFile,
		OutputDir: outputDir,
		Package:   pkgName,
		Verbose:   verbose,
	}

	if err := generator.Generate(config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
