package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// StringValidator validates string values
type StringValidator struct {
	minLen     *int
	maxLen     *int
	pattern    *regexp.Regexp
	email      bool
	custom     func(string) *Error
	required   bool
	defaultVal *string
	optional   bool
}

var _ Validator[string] = (*StringValidator)(nil)

// String creates a new string validator
func String() *StringValidator {
	return &StringValidator{}
}

// MinLen adds a minimum length validation rule
func (v *StringValidator) MinLen(length int) *StringValidator {
	v.minLen = &length
	return v
}

// MaxLen adds a maximum length validation rule
func (v *StringValidator) MaxLen(length int) *StringValidator {
	v.maxLen = &length
	return v
}

// Pattern adds a regular expression pattern validation rule
func (v *StringValidator) Pattern(pattern string) *StringValidator {
	v.pattern = regexp.MustCompile(pattern)
	return v
}

// Matches adds a regular expression pattern validation rule (alias for Pattern)
func (v *StringValidator) Matches(pattern string) *StringValidator {
	return v.Pattern(pattern)
}

// Default sets a default value to use if the string is empty
func (v *StringValidator) Default(val string) *StringValidator {
	v.defaultVal = &val
	return v
}

// Optional marks the field as optional (allows empty strings)
func (v *StringValidator) Optional() *StringValidator {
	v.optional = true
	return v
}

// Catch sets a fallback value to use if validation fails
func (v *StringValidator) Catch(val string) *StringValidator {
	// For now, we'll store the catch value but not implement the full logic
	// This would need more complex error handling
	return v
}

// Email adds an email validation rule
func (v *StringValidator) Email() *StringValidator {
	v.email = true
	return v
}

// Required adds a required field validation rule
func (v *StringValidator) Required() *StringValidator {
	v.required = true
	return v
}

// Custom adds a custom validation rule
func (v *StringValidator) Custom(fn func(string) *Error) *StringValidator {
	v.custom = fn
	return v
}

// Validate implements the Validator interface
func (v *StringValidator) Validate(value string) *Error {
	// Apply default if value is empty and default is set
	if v.defaultVal != nil && len(strings.TrimSpace(value)) == 0 {
		value = *v.defaultVal
	}

	// Check if required
	if v.required && len(strings.TrimSpace(value)) == 0 {
		return &Error{
			Code:    "required",
			Message: "field is required",
		}
	}

	// If optional and empty, skip validation
	if v.optional && len(strings.TrimSpace(value)) == 0 {
		return nil
	}

	if v.minLen != nil {
		if len(value) < *v.minLen {
			return &Error{
				Code:    "too_short",
				Message: fmt.Sprintf("must be at least %d characters", *v.minLen),
			}
		}
	}

	if v.maxLen != nil {
		if len(value) > *v.maxLen {
			return &Error{
				Code:    "too_long",
				Message: fmt.Sprintf("must be at most %d characters", *v.maxLen),
			}
		}
	}

	if v.pattern != nil {
		if !v.pattern.MatchString(value) {
			return &Error{
				Code:    "invalid_format",
				Message: "invalid format",
			}
		}
	}

	if v.email {
		if !strings.Contains(value, "@") || !strings.Contains(value, ".") {
			return &Error{
				Code:    "invalid_email",
				Message: "must be a valid email address",
			}
		}
	}

	if v.custom != nil {
		if err := v.custom(value); err != nil {
			return err
		}
	}

	return nil
}
