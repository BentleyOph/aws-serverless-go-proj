package validators

import (
	"regexp"
)

// IsEmail is a function that checks if the given string is a valid email address.
func IsEmail(email string) bool {
	// This regular expression is from https://emailregex.com/
	// It is a simple regular expression that checks if the given string is a valid email address.
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}