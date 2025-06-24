package main

import (
	"encoding/json"
	"fmt"

	"github.com/bm-197/tibeb/pkg/validate"
)

type User struct {
	Username string
	Email    string
	Age      int
}

func main() {
	// Create a validation schema for User
	schema := validate.Struct[User]().
		Field(func(u User) string { return u.Username }, validate.String().MinLen(3).MaxLen(30)).
		Field(func(u User) string { return u.Email }, validate.String().Email()).
		Field(func(u User) int { return u.Age }, validate.Int().Min(13))

	// Valid user
	validUser := User{
		Username: "johndoe",
		Email:    "john@example.com",
		Age:      25,
	}

	if errs := schema.Validate(validUser); errs.HasErrors() {
		printErrors("Valid user validation (should pass):", errs)
	} else {
		fmt.Println("Valid user passed validation!")
	}

	// Invalid user
	invalidUser := User{
		Username: "jo", // too short
		Email:    "not-an-email",
		Age:      10, // too young
	}

	if errs := schema.Validate(invalidUser); errs.HasErrors() {
		printErrors("Invalid user validation (should fail):", errs)
	} else {
		fmt.Println("Invalid user passed validation (unexpected)!")
	}
}

func printErrors(header string, errs *validate.Errors) {
	fmt.Println(header)
	errJSON, _ := json.MarshalIndent(errs.Get(), "", "  ")
	fmt.Println(string(errJSON))
}
