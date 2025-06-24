package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bm-197/tibeb/pkg/validate"
)

type Address struct {
	Street  string
	City    string
	Country string
	ZipCode string
}

type User struct {
	Username string
	Email    string
	Age      int
	Address  Address
	Role     string
}

func main() {
	// Create an address schema
	addressSchema := validate.Struct[Address]().
		Field(func(a Address) string { return a.Street }, validate.String().MinLen(5)).
		Field(func(a Address) string { return a.City }, validate.String().Required()).
		Field(func(a Address) string { return a.Country }, validate.String().Required()).
		Field(func(a Address) string { return a.ZipCode }, validate.Custom(validateZipCode))

	// Create a user schema with nested validation and composition
	userSchema := validate.Struct[User]().
		// Username validation with custom rule
		Field(func(u User) string { return u.Username }, validate.Custom(validateUsername)).
		// Email validation
		Field(func(u User) string { return u.Email }, validate.String().Email()).
		// Age validation with composition
		Field(func(u User) int { return u.Age },
			validate.OneOf(
				validate.Int().Min(13).Max(19), // Teenager
				validate.Int().Min(65),         // Senior
			),
		).
		// Nested address validation
		Field(func(u User) Address { return u.Address }, validate.Nested(addressSchema)).
		// Role validation with Not
		Field(func(u User) string { return u.Role },
			validate.AllOf(
				validate.String().Required(),
				validate.Not(validate.String().Matches("(?i)admin")), // Cannot be admin
			))

	// Valid teenage user
	validUser := User{
		Username: "john_doe_2024",
		Email:    "john@example.com",
		Age:      15,
		Address: Address{
			Street:  "123 Main Street",
			City:    "New York",
			Country: "USA",
			ZipCode: "12345",
		},
		Role: "user",
	}

	if errs := userSchema.Validate(validUser); errs.HasErrors() {
		printErrors("Valid user validation (should pass):", errs)
	} else {
		fmt.Println("Valid user passed validation!")
	}

	// Invalid user
	invalidUser := User{
		Username: "ad", // too short
		Email:    "not-an-email",
		Age:      30, // neither teenager nor senior
		Address: Address{
			Street:  "123", // too short
			City:    "",    // required
			Country: "",    // required
			ZipCode: "abc", // invalid format
		},
		Role: "admin", // not allowed
	}

	if errs := userSchema.Validate(invalidUser); errs.HasErrors() {
		printErrors("Invalid user validation (should fail):", errs)
	} else {
		fmt.Println("Invalid user passed validation (unexpected)!")
	}
}

// Custom validation function for username
func validateUsername(username string) *validate.Error {
	if len(username) < 8 {
		return &validate.Error{
			Code:    "username_too_short",
			Message: "username must be at least 8 characters",
		}
	}
	if !strings.Contains(username, "_") {
		return &validate.Error{
			Code:    "invalid_username_format",
			Message: "username must contain an underscore",
		}
	}
	return nil
}

// Custom validation function for zip code
func validateZipCode(zipCode string) *validate.Error {
	if len(zipCode) != 5 {
		return &validate.Error{
			Code:    "invalid_zipcode_length",
			Message: "zip code must be exactly 5 characters",
		}
	}
	for _, c := range zipCode {
		if c < '0' || c > '9' {
			return &validate.Error{
				Code:    "invalid_zipcode_format",
				Message: "zip code must contain only digits",
			}
		}
	}
	return nil
}

func printErrors(header string, errs *validate.Errors) {
	fmt.Println(header)
	errJSON, _ := json.MarshalIndent(errs.Get(), "", "  ")
	fmt.Println(string(errJSON))
}
