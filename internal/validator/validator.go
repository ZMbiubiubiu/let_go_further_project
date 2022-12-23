package validator

import "regexp"

// Declare a regular expression for sanity checking the format of email addresses(we'll use this later in the book).
// If you're interested, this regular expression pattern is taken from  https://html.spec.whatwg.org/#valid-e-mail-address.
// Note: if you're reading this if PDF or EPUB format and cannot seet he full pattern, please see the note further down the page.
var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Define a new Validator type which contains a map of validation errors.
type Validator struct {
	Errors map[string]string
}

// New is hepler which creates a new Validator instance with an empty errors map
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if the errors map doesn't contain any entry
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map
func (v *Validator) AddError(key, message string) {
	if _, isExists := v.Errors[key]; !isExists {
		v.Errors[key] = message
	}
}

// Check adds an error message to the map only if a validation check is not 'ok'
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// Generic function which returns true if a specific value is in a list.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// Mathces return true if a string value matches a specific regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Generic function which returns true if all values in a slice are unique
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
