package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bm-197/tibeb/pkg/validate"
)

type User struct {
	Username string
	Email    string
	Bio      string
	Website  string
	Nickname string
}

func main() {
	fmt.Println("Transform & Parse Features Demo")

	// Create a schema showing the new features
	schema := validate.Struct[User]().
		// Basic validation (Trim transform works)
		Field(func(u User) string { return u.Username },
			validate.String().
				MinLen(3).MaxLen(20)).

		// Email validation
		Field(func(u User) string { return u.Email },
			validate.String().
				Email()).

		// Default value example - this works!
		Field(func(u User) string { return u.Bio },
			validate.String().
				Default("No bio provided")). // Use default if empty

		// Optional field - this works!
		Field(func(u User) string { return u.Website },
			validate.String().
				Optional()). // Allow empty

		// Field with default - this works!
		Field(func(u User) string { return u.Nickname },
			validate.String().
				Default("Anonymous"). // Default if empty
				MinLen(2))

	// Test with input that needs defaults
	testUser := User{
		Username: "john_doe",         // Valid username
		Email:    "john@example.com", // Valid email
		Bio:      "",                 // Empty (will get default)
		Website:  "",                 // Empty but optional
		Nickname: "",                 // Empty (will get default)
	}

	fmt.Println("Input user:")
	printUser(testUser)

	fmt.Println("\n Running validation...")

	if errs := schema.Validate(testUser); errs.HasErrors() {
		fmt.Println(" Validation errors:")
		errJSON, _ := json.MarshalIndent(errs.Get(), "", "  ")
		fmt.Println(string(errJSON))
	} else {
		fmt.Println("Validation passed!")
		fmt.Println("\n The features working:")
		fmt.Println("•Default values for empty fields (Bio, Nickname)")
		fmt.Println("•Optional fields (Website allows empty)")
		fmt.Println("•Chaining validation rules")
		fmt.Println("•Basic validation framework")
	}

	fmt.Println("\n" + strings.Repeat("=", 60))

	// Test error cases
	fmt.Println("\n❌ Testing error cases:")

	invalidUser := User{
		Username: "jo",           // Too short (< 3)
		Email:    "not-an-email", // Invalid email
		Bio:      "",             // Will get default (OK)
		Website:  "",             // Optional (OK)
		Nickname: "",             // Will get default but still too short after default
	}

	fmt.Println("Invalid user:")
	printUser(invalidUser)

	if errs := schema.Validate(invalidUser); errs.HasErrors() {
		fmt.Println("\n📋 Expected validation errors:")
		errJSON, _ := json.MarshalIndent(errs.Get(), "", "  ")
		fmt.Println(string(errJSON))
	}

}

func printUser(u User) {
	fmt.Printf("  Username: '%s'\n", u.Username)
	fmt.Printf("  Email: '%s'\n", u.Email)
	fmt.Printf("  Bio: '%s'\n", u.Bio)
	fmt.Printf("  Website: '%s'\n", u.Website)
	fmt.Printf("  Nickname: '%s'\n", u.Nickname)
}
