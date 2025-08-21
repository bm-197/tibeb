package validate

import (
	"time"
)

// TimeValidator validates time.Time values
type TimeValidator struct {
	after    *time.Time
	before   *time.Time
	between  *[2]time.Time
	custom   func(time.Time) *Error
	required bool
}

var _ Validator[time.Time] = (*TimeValidator)(nil)

// Time creates a new time validator
func Time() *TimeValidator {
	return &TimeValidator{}
}

// After adds validation that time must be after the given time
func (v *TimeValidator) After(t time.Time) *TimeValidator {
	v.after = &t
	return v
}

// Before adds validation that time must be before the given time
func (v *TimeValidator) Before(t time.Time) *TimeValidator {
	v.before = &t
	return v
}

// Between adds validation that time must be between two times
func (v *TimeValidator) Between(start, end time.Time) *TimeValidator {
	v.between = &[2]time.Time{start, end}
	return v
}

// Custom adds a custom validation function
func (v *TimeValidator) Custom(fn func(time.Time) *Error) *TimeValidator {
	v.custom = fn
	return v
}

// Required marks the field as required
func (v *TimeValidator) Required() *TimeValidator {
	v.required = true
	return v
}

// Validate validates a time value
func (v *TimeValidator) Validate(value time.Time) *Error {
	// Check if required
	if v.required && value.IsZero() {
		return &Error{
			Field:   "",
			Code:    "required",
			Message: "field is required",
		}
	}

	// Skip validation for zero time if not required
	if !v.required && value.IsZero() {
		return nil
	}

	// Check after constraint
	if v.after != nil && !value.After(*v.after) {
		return &Error{
			Field:   "",
			Code:    "too_early",
			Message: "time must be after " + v.after.Format(time.RFC3339),
		}
	}

	// Check before constraint
	if v.before != nil && !value.Before(*v.before) {
		return &Error{
			Field:   "",
			Code:    "too_late",
			Message: "time must be before " + v.before.Format(time.RFC3339),
		}
	}

	// Check between constraint
	if v.between != nil {
		start, end := v.between[0], v.between[1]
		if value.Before(start) || value.After(end) {
			return &Error{
				Field:   "",
				Code:    "out_of_range",
				Message: "time must be between " + start.Format(time.RFC3339) + " and " + end.Format(time.RFC3339),
			}
		}
	}

	// Check custom validation
	if v.custom != nil {
		if err := v.custom(value); err != nil {
			return err
		}
	}

	return nil
}

// Common time validation helpers
func (v *TimeValidator) Today() *TimeValidator {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)
	return v.Between(start, end)
}

func (v *TimeValidator) Future() *TimeValidator {
	return v.After(time.Now())
}

func (v *TimeValidator) Past() *TimeValidator {
	return v.Before(time.Now())
}

func (v *TimeValidator) BusinessDay() *TimeValidator {
	return v.Custom(func(t time.Time) *Error {
		weekday := t.Weekday()
		if weekday == time.Saturday || weekday == time.Sunday {
			return &Error{
				Field:   "",
				Code:    "not_business_day",
				Message: "must be a business day (Monday-Friday)",
			}
		}
		return nil
	})
}
