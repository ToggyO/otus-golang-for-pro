package apperrors

import (
	"errors"
	"fmt"
)

var (
	ErrArgumentNil = errors.New("argument must not be nil value")
)

func ErrNotFound(object interface{}) error {
	return errors.New(fmt.Sprintf("not found: %v", object))
}
