package validate

// OneOfValidator checks if at least one validator passes
type OneOfValidator[T any] struct {
	validators []Validator[T]
}

// OneOf creates a new validator that passes if any of the given validators pass
func OneOf[T any](validators ...Validator[T]) Validator[T] {
	return &OneOfValidator[T]{
		validators: validators,
	}
}

// Validate implements the Validator interface
func (v *OneOfValidator[T]) Validate(value T) *Error {
	var lastError *Error
	for _, validator := range v.validators {
		if err := validator.Validate(value); err == nil {
			return nil
		} else {
			lastError = err
		}
	}
	return &Error{
		Code:    "no_match",
		Message: "value did not match any of the requirements",
		Field:   lastError.Field,
	}
}

// AllOfValidator checks if all validators pass
type AllOfValidator[T any] struct {
	validators []Validator[T]
}

// AllOf creates a new validator that passes if all of the given validators pass
func AllOf[T any](validators ...Validator[T]) Validator[T] {
	return &AllOfValidator[T]{
		validators: validators,
	}
}

// Validate implements the Validator interface
func (v *AllOfValidator[T]) Validate(value T) *Error {
	for _, validator := range v.validators {
		if err := validator.Validate(value); err != nil {
			return err
		}
	}
	return nil
}

// NotValidator inverts the result of another validator
type NotValidator[T any] struct {
	validator Validator[T]
}

// Not creates a new validator that passes if the given validator fails
func Not[T any](validator Validator[T]) Validator[T] {
	return &NotValidator[T]{
		validator: validator,
	}
}

// Validate implements the Validator interface
func (v *NotValidator[T]) Validate(value T) *Error {
	if err := v.validator.Validate(value); err == nil {
		return &Error{
			Code:    "invalid_match",
			Message: "value matched when it should not have",
		}
	}
	return nil
} 