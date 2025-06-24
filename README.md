# Tibeb

Tibeb is a type-safe, composable validation library for Go, inspired by libraries like Zod (JavaScript). It leverages Go generics and first-class functions to provide a clean, fluent API for defining validation rules.

## Features

- 🔒 Fully type-safe validation using Go generics
- 🧩 Composable validation rules
- 🎯 No reflection-based struct tags
- 📝 Clear, structured error messages
- 🚀 Simple, fluent API

## Installation

```bash
go get github.com/yourusername/tibeb
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/yourusername/tibeb/pkg/validate"
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
    Required()           // Non-empty
```

### Integer Validator
```go
validate.Int().
    Min(13).     // Minimum value
    Max(100).    // Maximum value
    Positive().  // Must be > 0
    Negative()   // Must be < 0
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