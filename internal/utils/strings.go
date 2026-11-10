package utils

import "unicode"

func MixCase(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		if i%2 == 0 {
			runes[i] = unicode.ToUpper(r)
		} else {
			runes[i] = unicode.ToLower(r)
		}
	}
	return string(runes)
}
