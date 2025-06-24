package validate

// NestedValidator provides validation for nested structs
type NestedValidator[T any] struct {
	schema *Schema[T]
}

// Nested creates a new nested struct validator
func Nested[T any](schema *Schema[T]) Validator[T] {
	return &NestedValidator[T]{
		schema: schema,
	}
}

// Validate implements the Validator interface
func (v *NestedValidator[T]) Validate(value T) *Error {
	if errs := v.schema.Validate(value); errs.HasErrors() {
		// Return the first error with the proper field path
		firstErr := errs.Get()[0]
		return &Error{
			Code:    firstErr.Code,
			Message: firstErr.Message,
			Field:   firstErr.Field,
		}
	}
	return nil
}
