package main

import (
	"encoding/json"
	"fmt"

	"github.com/bm-197/tibeb/examples/codegen/models"
	"github.com/bm-197/tibeb/pkg/validate"
)

func main() {
	// Valid user
	validUser := models.User{
		Username: "johndoe",
		Email:    "john@example.com",
		Age:      25,
	}

	if errs := models.ValidateUser(validUser); errs.HasErrors() {
		printErrors("Valid user validation (should pass):", errs)
	} else {
		fmt.Println("Valid user passed validation!")
	}

	// Invalid user
	invalidUser := models.User{
		Username: "jo", // too short
		Email:    "not-an-email",
		Age:      10, // too young
	}

	if errs := models.ValidateUser(invalidUser); errs.HasErrors() {
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
