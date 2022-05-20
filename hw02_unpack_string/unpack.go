package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var nextChar rune
	var mustBeEscaped bool
	var escapeCharRune rune = 92
	var result strings.Builder

	runes := []rune(str)
	length := len(runes)

	for i := 0; i < length; i++ {
		currentChar := runes[i]
		isCurrentDigit := unicode.IsDigit(currentChar)

		if i < length-1 {
			nextChar = runes[i+1]
		} else {
			nextChar = 0
		}

		isStartsWithDigit := isCurrentDigit && i == 0
		isInvalidDigit := isCurrentDigit && unicode.IsDigit(nextChar) && !mustBeEscaped
		isInvalidEscape := mustBeEscaped && !(currentChar == escapeCharRune || isCurrentDigit)

		if isStartsWithDigit || isInvalidDigit || isInvalidEscape {
			return "", ErrInvalidString
		}

		if mustBeEscaped || unicode.IsLetter(currentChar) || unicode.IsSymbol(currentChar) {
			mustBeEscaped = false

			if unicode.IsDigit(nextChar) {
				if err := unpackLetter(currentChar, nextChar, &result); err != nil {
					return "", err
				}
				continue
			}

			result.WriteRune(currentChar)
			continue
		}

		if currentChar == escapeCharRune {
			mustBeEscaped = true
		}
	}

	return result.String(), nil
}

func unpackLetter(char rune, numRune rune, result *strings.Builder) error {
	count, err := strconv.Atoi(string(numRune))
	if err != nil {
		return ErrInvalidString
	}

	for j := 0; j < count; j++ {
		result.WriteRune(char)
	}

	return nil
}
