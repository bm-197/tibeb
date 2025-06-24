package validate

import "reflect"

// Struct creates a new schema for validating structs of type T
func Struct[T any]() *Schema[T] {
	return &Schema[T]{
		rules: make([]FieldRule[T], 0),
	}
}

// Field adds a field validation rule to the schema
func (s *Schema[T]) Field(selector interface{}, validator interface{}) *Schema[T] {
	// Get the field name from the selector function using reflection
	t := reflect.TypeOf((*T)(nil)).Elem()
	selectorVal := reflect.ValueOf(selector)

	if selectorVal.Kind() != reflect.Func {
		panic("selector must be a function")
	}

	// Extract field name from the selector
	fieldName := ""
	if t.Kind() == reflect.Struct {
		// Create a zero value of type T
		var zero T
		zeroVal := reflect.ValueOf(zero)
		result := selectorVal.Call([]reflect.Value{zeroVal})[0]
		selectorType := result.Type()

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Type == selectorType {
				fieldName = field.Name
				break
			}
		}
	}

	// Create a wrapper that converts the field value to any
	wrapper := func(t T) any {
		result := selectorVal.Call([]reflect.Value{reflect.ValueOf(t)})[0]
		return result.Interface()
	}

	// Create a wrapper that converts the validator to handle any
	validatorVal := reflect.ValueOf(validator)
	validateMethod := validatorVal.MethodByName("Validate")
	if !validateMethod.IsValid() {
		panic("validator must implement Validate method")
	}

	validatorWrapper := ValidatorFunc[any](func(value any) *Error {
		result := validateMethod.Call([]reflect.Value{reflect.ValueOf(value)})
		if len(result) != 1 {
			panic("Validate method must return exactly one value")
		}
		if result[0].IsNil() {
			return nil
		}
		return result[0].Interface().(*Error)
	})

	s.rules = append(s.rules, FieldRule[T]{
		selector: wrapper,
		rule:     validatorWrapper,
		field:    fieldName,
	})

	return s
}

// ValidatorFunc is a helper type that allows functions to implement Validator
type ValidatorFunc[T any] func(T) *Error

// Validate implements the Validator interface
func (f ValidatorFunc[T]) Validate(value T) *Error {
	return f(value)
}

// TypedField is a helper function to create a typed field validator
func TypedField[T, F any](selector func(T) F, rule Validator[F]) (func(T) any, Validator[any]) {
	wrapper := func(t T) any {
		return selector(t)
	}

	validatorWrapper := ValidatorFunc[any](func(value any) *Error {
		if v, ok := value.(F); ok {
			return rule.Validate(v)
		}
		return &Error{
			Code:    "invalid_type",
			Message: "invalid field type",
		}
	})

	return wrapper, validatorWrapper
}
