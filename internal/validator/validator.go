// Filename: test2/internal/data/validator.go
package validator

import (
	"net/url"
	"regexp"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

// creates a new instance
func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// checks the errors map for any entries
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// Checks if an element can be found in provied list of elements
func In(element string, list ...string) bool {
	for i := range list {
		if element == list[i] {
			return true
		}
	}
	return false
}

// Matches returns true if a string value matches a specific regex pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// checks if a string value is valid web URL
func ValidWebsite(website string) bool {
	_, err := url.ParseRequestURI(website)
	return err == nil
}

// adds an error entry into the errors map
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// performs the valiudation checks and class the adderror method in turn if error entry needs to be added
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// checks if there are no repeating values in the slice
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)
	for _, value := range values {
		uniqueValues[value] = true
	}
	return len(values) == len(uniqueValues)
}
