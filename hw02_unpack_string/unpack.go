package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	runes := []rune(s)
	slash, _ := utf8.DecodeRune([]byte("\\"))

	var (
		builder strings.Builder
		escape  bool
		rp      int
	)

	for rp < len(runes) {
		isInvalidEscapedChar := escape && runes[rp] != slash && !unicode.IsDigit(runes[rp])
		isUnescapedDigit := unicode.IsDigit(runes[rp]) && !escape

		switch {
		case isInvalidEscapedChar || isUnescapedDigit:
			return "", ErrInvalidString
		case runes[rp] == slash && !escape:
			escape = true
			rp++
		case rp+1 < len(runes) && unicode.IsDigit(runes[rp+1]):
			nRune := int(runes[rp+1] - '0')
			repeatedRune := strings.Repeat(string(runes[rp]), nRune)
			builder.WriteString(repeatedRune)
			escape = false
			rp += 2
		default:
			builder.WriteRune(runes[rp])
			escape = false
			rp++
		}
	}

	if escape {
		return "", ErrInvalidString
	}

	return builder.String(), nil
}
