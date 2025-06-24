package validate

import "fmt"

// IntValidator provides validation rules for integer values
type IntValidator struct {
	min      *int
	max      *int
	positive bool
	negative bool
}

var _ Validator[int] = (*IntValidator)(nil)

// Int creates a new integer validator
func Int() *IntValidator {
	return &IntValidator{}
}

// Min adds a minimum value validation rule
func (v *IntValidator) Min(value int) *IntValidator {
	v.min = &value
	return v
}

// Max adds a maximum value validation rule
func (v *IntValidator) Max(value int) *IntValidator {
	v.max = &value
	return v
}

// Positive requires the value to be positive (> 0)
func (v *IntValidator) Positive() *IntValidator {
	v.positive = true
	return v
}

// Negative requires the value to be negative (< 0)
func (v *IntValidator) Negative() *IntValidator {
	v.negative = true
	return v
}

// Validate implements the Validator[int] interface
func (v *IntValidator) Validate(value int) *Error {
	if v.min != nil && value < *v.min {
		return &Error{
			Code:    "too_small",
			Message: fmt.Sprintf("value must be at least %d", *v.min),
		}
	}

	if v.max != nil && value > *v.max {
		return &Error{
			Code:    "too_large",
			Message: fmt.Sprintf("value must be at most %d", *v.max),
		}
	}

	if v.positive && value <= 0 {
		return &Error{
			Code:    "not_positive",
			Message: "value must be positive",
		}
	}

	if v.negative && value >= 0 {
		return &Error{
			Code:    "not_negative",
			Message: "value must be negative",
		}
	}

	return nil
}
