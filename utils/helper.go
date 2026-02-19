package utils

import (
	"crypto/rand"
	"math/big"
	"regexp"
	"strings"
	"unicode"
)

func TrimSpaces(s string) string {
	return strings.TrimSpace(s)
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func SanitizeEmail(email string) string {
	return ToLower(TrimSpaces(email))
}

func SanitizeUsername(username string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_\-\.]+`)
	return re.ReplaceAllString(TrimSpaces(username), "")
}

func SanitizeString(s string) string {
	s = TrimSpaces(s)
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}

func IsAlphaNumeric(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(strings.ToLower(email))
}

func DefaultString(s, def string) string {
	if IsEmpty(s) {
		return def
	}
	return s
}

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func GenerateRandomCode(length int) (string, error) {
	const digits = "0123456789"
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		ret[i] = digits[num.Int64()]
	}
	return string(ret), nil
}
