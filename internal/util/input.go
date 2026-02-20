package util

import (
	"countryinfo/internal/fp"
)

func IsAsciiChar(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func IsTwoLetterCountryCode(countryCode string) bool {
	return len(countryCode) == 2 && fp.ForAll([]rune(countryCode), func(r rune) bool {
		return IsAsciiChar(r)
	})
}
