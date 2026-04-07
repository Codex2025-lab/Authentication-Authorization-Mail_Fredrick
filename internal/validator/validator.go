package validator

import (
	"regexp"
)

type Validator struct {
	Errors map[string]string
}

// NewValidator returns a new Validator
func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Check adds an error if condition is false
func (v *Validator) Check(ok bool, field, message string) {
	if !ok {
		v.Errors[field] = message
	}
}

// IsEmpty returns true if no errors
func (v *Validator) IsEmpty() bool {
	return len(v.Errors) == 0
}

// AddError adds a single error
func (v *Validator) AddError(field, message string) {
	v.Errors[field] = message
}

// Matches checks if value matches regex
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// EmailRX is a simple regex for emails
var EmailRX = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)