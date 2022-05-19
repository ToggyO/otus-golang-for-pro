package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

var escapeCharsCollection = []rune{'a', 'b', 't', 'n', 'f', 'r', 'v'}

func Unpack(str string) (string, error) {
	// Standard escape char literals collection
	sb := strings.Builder{}

	runes := []rune(str)
	end := len(runes)

	var escapeRune rune = 92 // ASCII slash serial number
	var nextChar rune
	// Indicates whether the current character escaped
	var mustBeEscaped bool

	// Indicates whether the previous char was a `slash`
	// Helps to handle cases, when a slash char is escaped
	var isPrevRuneSlash bool

	for i := 0; i < end; i++ {
		char := runes[i]

		if i < end-1 {
			nextChar = runes[i+1]
		}

		// Handle raw string literal escape chars
		if char == escapeRune {
			if isPrevRuneSlash {
				goto Main
			}

			if containsRune(escapeCharsCollection, nextChar) {
				return errorResponse()
			}

			isPrevRuneSlash = true
			mustBeEscaped = true
			continue
		}

	Main:
		skipIteration, err := handleChar(char, nextChar, i, end, &mustBeEscaped, &isPrevRuneSlash, &sb)
		if err != nil {
			return "", err
		}
		if skipIteration {
			continue
		}

	}

	return sb.String(), nil
}

func handleChar(currentChar rune,
	nextChar rune,
	loopIterationNumber int,
	loopIterationEnd int,
	mustBeEscaped *bool,
	isPrevRuneSlash *bool,
	sb *strings.Builder,
) (skipIteration bool, err error) {
	currentCharIsDigit := !*mustBeEscaped && unicode.IsDigit(currentChar)
	nextCharIsDigit := (loopIterationNumber < loopIterationEnd-1) && unicode.IsDigit(nextChar)

	*isPrevRuneSlash = false
	*mustBeEscaped = false

	// Handles case, when unpacked chars count is greater than 9
	if (currentCharIsDigit && loopIterationNumber == 0) || (currentCharIsDigit && nextCharIsDigit) {
		return false, ErrInvalidString
	}

	// Handles case, when current and next chars is letters
	if !currentCharIsDigit && !nextCharIsDigit {
		sb.WriteRune(currentChar)
		return true, nil
	}

	// Handles case, when current char is letters, but next is a digit
	if !currentCharIsDigit && nextCharIsDigit {
		num, err := strconv.Atoi(string(nextChar))
		if err != nil {
			return false, ErrInvalidString
		}

		for j := 0; j < num; j++ {
			sb.WriteRune(currentChar)
		}
	}

	return false, nil
}

func containsRune(s []rune, e rune) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func errorResponse() (string, error) {
	return "", ErrInvalidString
}
