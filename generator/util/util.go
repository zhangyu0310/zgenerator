package util

import "unicode"

// FirstToUpper convert first letter to upper
func FirstToUpper(str string) string {
	runeStr := []rune(str)
	for i := 0; i < len(runeStr); i++ {
		if unicode.IsLetter(runeStr[i]) {
			if unicode.IsLower(runeStr[i]) {
				runeStr[i] = unicode.ToUpper(runeStr[i])
			}
			break
		}
	}
	return string(runeStr)
}
