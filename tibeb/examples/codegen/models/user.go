package models

import "github.com/bm-197/tibeb/pkg/validate"

type User struct {
	Username string
	Email    string
	Age      int
}

// Schema is the validation schema for User
var Schema = validate.Struct[User]().
	Field(func(u User) string { return u.Username }, validate.String().MinLen(3).MaxLen(30)).
	Field(func(u User) string { return u.Email }, validate.String().Email()).
	Field(func(u User) int { return u.Age }, validate.Int().Min(13))
