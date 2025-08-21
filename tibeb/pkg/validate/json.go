package validate

import (
	"encoding/json"
)

// JSONValidator validates that a value is valid JSON
type JSONValidator struct {
	custom func(interface{}) *Error
}

var _ Validator[interface{}] = (*JSONValidator)(nil)

// JSON creates a new JSON validator
func JSON() *JSONValidator {
	return &JSONValidator{}
}

// Custom adds a custom validation function for the parsed JSON
func (v *JSONValidator) Custom(fn func(interface{}) *Error) *JSONValidator {
	v.custom = fn
	return v
}

// Validate validates that the value is valid JSON
func (v *JSONValidator) Validate(value interface{}) *Error {
	// If it's a string, try to parse it as JSON
	if str, ok := value.(string); ok {
		var temp interface{}
		if err := json.Unmarshal([]byte(str), &temp); err != nil {
			return &Error{
				Field:   "",
				Code:    "invalid_json",
				Message: "invalid JSON format: " + err.Error(),
			}
		}
		value = temp
	}

	// Run custom validation on the parsed JSON
	if v.custom != nil {
		if err := v.custom(value); err != nil {
			return err
		}
	}

	return nil
}

// Convenience methods for common JSON validations
func (v *JSONValidator) Object() *JSONValidator {
	return v.Custom(func(val interface{}) *Error {
		if _, ok := val.(map[string]interface{}); !ok {
			return &Error{
				Field:   "",
				Code:    "not_object",
				Message: "must be a JSON object",
			}
		}
		return nil
	})
}

func (v *JSONValidator) Array() *JSONValidator {
	return v.Custom(func(val interface{}) *Error {
		if _, ok := val.([]interface{}); !ok {
			return &Error{
				Field:   "",
				Code:    "not_array",
				Message: "must be a JSON array",
			}
		}
		return nil
	})
}
