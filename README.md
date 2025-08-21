# Tibeb

Tibeb is a type-safe, composable validation library for Go, inspired by libraries like Zod (JavaScript). It leverages Go generics and first-class functions to provide a clean, fluent API for defining validation rules.

## Features

- 🔒 **Fully type-safe** validation using Go generics
- 🧩 **Composable** validation rules with fluent chaining
- 🎯 **No reflection** - uses function selectors instead of struct tags
- 📝 **Clear, structured** error messages with JSON support
- 🚀 **Simple, fluent API** inspired by Zod
- 🔄 **Transform & Parse** - Clean and convert data during validation
- ⚡ **Code Generation** - Generate optimized validators for production
- 🎨 **Default Values** - Automatic fallbacks for empty fields
- 📊 **Advanced Types** - Time, JSON, and custom validation support

## Installation

```bash
go get github.com/bm-197/tibeb
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/bm-197/tibeb/tibeb/pkg/validate"
)

type User struct {
    Username string
    Email    string
    Age      int
}

func main() {
    // Create a validation schema
    schema := validate.Struct[User]().
        Field(func(u User) string { return u.Username }, validate.String().MinLen(3).MaxLen(30)).
        Field(func(u User) string { return u.Email }, validate.String().Email()).
        Field(func(u User) int { return u.Age }, validate.Int().Min(13))

    // Validate a user
    user := User{
        Username: "jo",          // too short
        Email:    "not-an-email",
        Age:      10,           // too young
    }

    if errs := schema.Validate(user); errs.HasErrors() {
        // Handle validation errors
        for _, err := range errs.Get() {
            fmt.Printf("Field %s: %s\n", err.Field, err.Message)
        }
    }
}
```

## Available Validators

### String Validator
```go
validate.String().
    MinLen(3).           // Minimum length
    MaxLen(30).          // Maximum length
    Email().             // Email format
    Matches("^[a-z]+$"). // Regex pattern
    Required().          // Non-empty
    Optional().          // Allow empty
    Default("fallback"). // Default value if empty
    Trim().              // Remove whitespace
    Lowercase().         // Convert to lowercase
    Uppercase()          // Convert to uppercase
```

### Integer Validator
```go
validate.Int().
    Min(13).     // Minimum value
    Max(100).    // Maximum value
    Positive().  // Must be > 0
    Negative()   // Must be < 0
```

### Time Validator
```go
validate.Time().
    After(time.Now()).           // Must be in future
    Before(deadline).            // Must be before date
    Between(start, end).         // Must be in range
    BusinessDay().               // Monday-Friday only
    Future().                    // Must be in future
    Past()                       // Must be in past
```

### JSON Validator
```go
validate.JSON().
    Object().    // Must be JSON object
    Array()      // Must be JSON array
```

### Transform & Parse
```go
// Transform data during validation
validate.String().
    Transform(strings.ToLower).
    Transform(strings.TrimSpace).
    MinLen(3)

// Parse strings to other types
validate.String().ParseInt()                    // string → int
validate.String().ParseTime("2006-01-02")      // string → time.Time
validate.String().ParseJSON(target)            // string → JSON
```

## Code Generation (v0.3)

Generate optimized, zero-reflection validators for production:

```bash
# Install the CLI
go install github.com/yourusername/tibeb/cmd/tibeb

# Generate validators
tibeb gen -file=models/user.go -out=generated/
```

Input schema:
```go
// models/user.go
var Schema = validate.Struct[User]().
    Field(func(u User) string { return u.Username }, validate.String().MinLen(3)).
    Field(func(u User) string { return u.Email }, validate.String().Email())
```

Generated code:
```go
// generated/user_validator.go
func ValidateUser(v User) *validate.Errors {
    return UserSchema.Validate(v)
}

var UserSchema = validate.Struct[User]().
    Field(func(v User) string { return v.Username }, validate.String().MinLen(3)).
    Field(func(v User) string { return v.Email }, validate.String().Email())
```

## Advanced Examples

### Transform & Default Values
```go
schema := validate.Struct[User]().
    Field(func(u User) string { return u.Username },
        validate.String().
            Trim().                    // Remove whitespace
            Lowercase().               // Convert to lowercase
            Default("anonymous").      // Use default if empty
            MinLen(3)).
    Field(func(u User) string { return u.Bio },
        validate.String().
            Optional().                // Allow empty
            Default("No bio provided"))
```

### Custom Validation
```go
validate.String().Custom(func(s string) *validate.Error {
    if !isValidUsername(s) {
        return &validate.Error{
            Code:    "invalid_username",
            Message: "username contains invalid characters",
        }
    }
    return nil
})
```

## Error Handling

Validation errors are structured and can be easily converted to JSON:

```json
[
  {
    "field": "Username",
    "code": "too_short",
    "message": "length must be at least 3 characters"
  },
  {
    "field": "Email",
    "code": "invalid_email",
    "message": "invalid email address"
  },
  {
    "field": "Age",
    "code": "too_small",
    "message": "value must be at least 13"
  }
]
```
