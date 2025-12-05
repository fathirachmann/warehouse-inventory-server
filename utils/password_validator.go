package utils

import (
	"unicode"
)

// ValidatePassword mengecek kompleksitas password berdasarkan aturan:
// minimal 8 karakter, mengandung huruf, angka, dan karakter spesial
func ValidatePassword(pass string) bool {
	// 1. Cek panjang minimal
	if len(pass) < 8 {
		return false
	}

	var (
		hasLetter  bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range pass {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
		case unicode.IsNumber(char):
			hasNumber = true
		case isSpecialChar(char):
			hasSpecial = true
		default:
			// Jika ketemu karakter di luar A-Z, 0-9, atau spesial char yg ditentukan
			// Return false (karena regex kamu membatasi character set-nya)
			return false
		}
	}

	return hasLetter && hasNumber && hasSpecial
}

func isSpecialChar(c rune) bool {
	specialChars := "@$!%*#?&"
	for _, sc := range specialChars {
		if c == sc {
			return true
		}
	}
	return false
}
