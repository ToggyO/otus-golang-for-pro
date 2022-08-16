package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}

	for _, valErr := range v {
		sb.WriteString(fmt.Sprintf("ValidationError: %s\n", valErr.Err))
	}

	return sb.String()
}

var (
	ErrNotStruct                = errors.New("value is not a struct type")
	ErrUnsupportedType          = errors.New("unsupported field type: %s")
	ErrUnsupportedValidatorType = errors.New("unsupported validator type: %s")

	ErrValidationStringLength = errors.New("invalid string length of field `%s`, must be equal: %d")
	ErrValidationStringRegexp = errors.New("field `%s` doesnt't match provided regexp: %v")
	ErrValidationMin          = errors.New("field `%s` must be greater then: %v")
	ErrValidationMax          = errors.New("field `%s` must be less then: %v")
	ErrValidationIn           = errors.New("field `%s` is not in: {%v}")
)
