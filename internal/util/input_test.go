package util

import "testing"

func TestIsAsciiChar(t *testing.T) {
	tests := []struct {
		name string
		args rune
		want bool
	}{
		{"Unicode is not ascii", 'å', false},
		{"a is ascii", 'a', true},
		{"z is ascii", 'z', true},
		{"A is ascii", 'A', true},
		{"Z is ascii", 'Z', true},
		{"a-1 is not ascii char", 'a' - 1, false},
		{"z+1 is not ascii char", 'z' + 1, false},
		{"A-1 is not ascii char", 'A' - 1, false},
		{"Z+1 is not ascii char", 'Z' + 1, false},
		{"digit is not ascii char", '5', false},
		{"space is not ascii char", ' ', false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAsciiChar(tt.args); got != tt.want {
				t.Errorf("IsAsciiChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTwoLetterCountryCode(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{"no is valid country code", "no", true},
		{"nok is invalid country code", "nok", false},
		{"us is valid country code", "us", true},
		{"c is invlaid country code", "c", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTwoLetterCountryCode(tt.args); got != tt.want {
				t.Errorf("IsTwoLetterCountryCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
