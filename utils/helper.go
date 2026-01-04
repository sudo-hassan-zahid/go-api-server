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
