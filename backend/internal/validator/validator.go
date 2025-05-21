package validator

import (
	"regexp"
	"slices"
)

// EmailRX is a regex pattern for validating email addresses.
var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Validator maintains a map of validation errors.
type Validator struct {
	// Errors stores validation error messages keyed by field name
	Errors map[string]string
}

// New creates and returns a new Validator instance with an initialized Errors map.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if there are no validation errors, false otherwise.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the Errors map if the key doesn't already exists.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check evaluates a condition and adds an error message if the condition is false.
// This provides a convenient way to add validation errors conditionally.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In checks if a string value is contained within a list of allowed values.
// Returns true if the value is found in the list, false otherwise.
func In(value string, list ...string) bool {
	return slices.Contains(list, value)
}

// Matches checks if a string value matches a regex pattern.
// Returns true if there's a match, false otherwise.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique checks if all string values in a slice are distinct.
// Returns true if all values are unique, false if there are duplicates.
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
