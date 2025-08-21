package validate

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// TransformFunc represents a transformation function
type TransformFunc[T any] func(T) T

// ParseFunc represents a parsing function that can fail
type ParseFunc[T, U any] func(T) (U, error)

// TransformValidator wraps another validator with transformations
type TransformValidator[T any] struct {
	validator  Validator[T]
	transforms []TransformFunc[T]
	defaultVal *T
	catchVal   *T
}

var _ Validator[string] = (*TransformValidator[string])(nil)

// Transform creates a new transform validator from an existing validator
func (v *StringValidator) Transform(fn func(string) string) *TransformValidator[string] {
	return &TransformValidator[string]{
		validator:  v,
		transforms: []TransformFunc[string]{fn},
	}
}

// Transform creates a new transform validator from an existing validator
func (v *IntValidator) Transform(fn func(int) int) *TransformValidator[int] {
	return &TransformValidator[int]{
		validator:  v,
		transforms: []TransformFunc[int]{fn},
	}
}

// Pipe adds another transformation to the chain
func (v *TransformValidator[T]) Pipe(fn TransformFunc[T]) *TransformValidator[T] {
	v.transforms = append(v.transforms, fn)
	return v
}

// Default sets a default value to use if the input is zero/empty
func (v *TransformValidator[T]) Default(val T) *TransformValidator[T] {
	v.defaultVal = &val
	return v
}

// Catch sets a fallback value to use if validation fails
func (v *TransformValidator[T]) Catch(val T) *TransformValidator[T] {
	v.catchVal = &val
	return v
}

// Validate applies transformations then validates
func (v *TransformValidator[T]) Validate(value T) *Error {
	if v.defaultVal != nil && isZeroValue(value) {
		value = *v.defaultVal
	}

	// Apply all transformations in order
	for _, transform := range v.transforms {
		value = transform(value)
	}

	// Validate the transformed value
	if err := v.validator.Validate(value); err != nil {
		if v.catchVal != nil {
			return v.validator.Validate(*v.catchVal)
		}
		return err
	}

	return nil
}

// Common string transformations
func (v *StringValidator) Trim() *TransformValidator[string] {
	return v.Transform(strings.TrimSpace)
}

func (v *StringValidator) Lowercase() *TransformValidator[string] {
	return v.Transform(strings.ToLower)
}

func (v *StringValidator) Uppercase() *TransformValidator[string] {
	return v.Transform(strings.ToUpper)
}

// Note: Chaining transforms would require more complex generic constraints
// For now, use single transforms from StringValidator directly

// ParseValidator handles parsing from one type to another
type ParseValidator[T, U any] struct {
	parseFunc ParseFunc[T, U]
	validator Validator[U]
}

// Parse creates a new parse validator
func Parse[T, U any](parseFunc ParseFunc[T, U], validator Validator[U]) *ParseValidator[T, U] {
	return &ParseValidator[T, U]{
		parseFunc: parseFunc,
		validator: validator,
	}
}

// Common parse functions for strings
func (v *StringValidator) ParseInt() *ParseValidator[string, int] {
	return Parse(func(s string) (int, error) {
		return strconv.Atoi(s)
	}, Int())
}

func (v *StringValidator) ParseTime(layout string) *ParseValidator[string, time.Time] {
	return Parse(func(s string) (time.Time, error) {
		return time.Parse(layout, s)
	}, &TimeValidator{})
}

func (v *StringValidator) ParseJSON(target interface{}) *ParseValidator[string, interface{}] {
	return Parse(func(s string) (interface{}, error) {
		err := json.Unmarshal([]byte(s), target)
		return target, err
	}, &JSONValidator{})
}

// Validate for ParseValidator
func (v *ParseValidator[T, U]) Validate(value T) *Error {
	parsed, err := v.parseFunc(value)
	if err != nil {
		return &Error{
			Field:   "",
			Code:    "parse_error",
			Message: "failed to parse value: " + err.Error(),
		}
	}

	return v.validator.Validate(parsed)
}

func isZeroValue[T any](value T) bool {
	var zero T
	return any(value) == any(zero)
}
