package validator

import (
	"reflect"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

func (v *Validator) SetValidator(validator *Validator) {
	v.FieldErrors = validator.FieldErrors
	v.NonFieldErrors = validator.NonFieldErrors
}

func New(model any) Validator {
	return Validator{
		FieldErrors:    initErrors(model),
		NonFieldErrors: []string{},
	}
}

func initErrors(model any) map[string]string {
	errors := make(map[string]string)

	// Use reflection to iterate over the fields of the model
	val := reflect.ValueOf(model)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		for i := 0; i < val.NumField(); i++ {
			field := val.Type().Field(i)
			errors[field.Name] = ""
		}
	}

	return errors
}

func (v *Validator) Valid() bool {
	// Check if the FieldErrors map is nil or empty
	if v.FieldErrors == nil {
		return true
	}

	// Check if there are any errors in the FieldErrors map
	for _, err := range v.FieldErrors {
		if err != "" {
			return false
		}
	}

	return true
}

func (v *Validator) AddFieldError(field, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	// Check if the field already has an error message
	if _, exists := v.FieldErrors[field]; !exists || v.FieldErrors[field] == "" {
		v.FieldErrors[field] = message
	}
}

// Create an AddNonFieldError() helper for adding error messages to the new
// NonFieldErrors slice.
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

func MinChars(value string, min int) bool {
	return utf8.RuneCountInString(value) >= min
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches() returns true if a value matches a provided compiled regular expression.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
