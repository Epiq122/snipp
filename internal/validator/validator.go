package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

// contains a map of validation error messages
type Validator struct {
	FieldErrors map[string]string
}

// returns true if the fieldErrors map is empty
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// adds an error to the map as long as no entry already exists for the field
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// adds error if validation check it not ok
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// adds an error to the map if the field is blank
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// adds an error to the map if the field is too long
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// returns true if the value is one of the allowed values
func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
