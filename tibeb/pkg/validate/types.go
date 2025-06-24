package validate

// Error represents a validation error
type Error struct {
	Field   string `json:"field,omitempty"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Errors represents a collection of validation errors
type Errors struct {
	errors []*Error
}

// Add adds a validation error to the collection
func (e *Errors) Add(err *Error) {
	if e.errors == nil {
		e.errors = make([]*Error, 0)
	}
	e.errors = append(e.errors, err)
}

// HasErrors returns true if there are any validation errors
func (e *Errors) HasErrors() bool {
	return len(e.errors) > 0
}

// Get returns all validation errors
func (e *Errors) Get() []*Error {
	return e.errors
}

// Validator is the interface for all validators
type Validator[T any] interface {
	Validate(value T) *Error
}

// Schema represents a validation schema for a struct
type Schema[T any] struct {
	rules []FieldRule[T]
}

// FieldRule represents a validation rule for a struct field
type FieldRule[T any] struct {
	selector func(T) any
	rule     Validator[any]
	field    string
}

// Validate runs all validators in the schema and returns any errors
func (s *Schema[T]) Validate(value T) *Errors {
	errors := &Errors{}
	for _, rule := range s.rules {
		if err := rule.rule.Validate(rule.selector(value)); err != nil {
			err.Field = rule.field
			errors.Add(err)
		}
	}
	return errors
}
