package utils

import "strings"

var fancyMap = map[rune]rune{
	'a': 'ᴀ', 'b': 'ʙ', 'c': 'ᴄ', 'd': 'ᴅ', 'e': 'ᴇ',
	'f': 'ғ', 'g': 'ɢ', 'h': 'ʜ', 'i': 'ɪ', 'j': 'ᴊ',
	'k': 'ᴋ', 'l': 'ʟ', 'm': 'ᴍ', 'n': 'ɴ', 'o': 'ᴏ',
	'p': 'ᴘ', 'q': 'ǫ', 'r': 'ʀ', 's': 's', 't': 'ᴛ',
	'u': 'ᴜ', 'v': 'ᴠ', 'w': 'ᴡ', 'x': 'x', 'y': 'ʏ', 'z': 'ᴢ',
}

func FancyText(input string) string {
	var out strings.Builder
	for _, r := range input {
		if r >= 'a' && r <= 'z' {
			if fr, ok := fancyMap[r]; ok {
				out.WriteRune(fr)
			} else {
				out.WriteRune(r)
			}
		} else {
			out.WriteRune(r)
		}
	}
	return out.String()
}
