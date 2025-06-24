package validate

// CustomValidator allows creating custom validation rules
type CustomValidator[T any] struct {
	validate func(T) *Error
}

// Custom creates a new custom validator
func Custom[T any](validate func(T) *Error) Validator[T] {
	return &CustomValidator[T]{
		validate: validate,
	}
}

// Validate implements the Validator interface
func (v *CustomValidator[T]) Validate(value T) *Error {
	return v.validate(value)
}
