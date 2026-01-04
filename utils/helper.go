package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// TrimSpaces removes leading and trailing spaces from a string
func TrimSpaces(s string) string {
	return strings.TrimSpace(s)
}

// ToLower converts a string to lowercase
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper converts a string to uppercase
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// SanitizeEmail trims spaces and lowercases the email
func SanitizeEmail(email string) string {
	return ToLower(TrimSpaces(email))
}

// SanitizeUsername removes invalid characters and trims the username
func SanitizeUsername(username string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_\-\.]+`)
	return re.ReplaceAllString(TrimSpaces(username), "")
}

// SanitizeString trims spaces and removes control characters
func SanitizeString(s string) string {
	s = TrimSpaces(s)
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}

// IsAlphaNumeric checks if a string contains only letters and numbers
func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

// IsValidEmail rudimentary email check
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(strings.ToLower(email))
}

// ------------------------ Misc ------------------------

// DefaultString returns the default if the string is empty
func DefaultString(s, def string) string {
	if IsEmpty(s) {
		return def
	}
	return s
}
