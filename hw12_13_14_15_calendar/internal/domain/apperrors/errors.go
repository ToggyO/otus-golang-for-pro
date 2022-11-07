package apperrors

import (
	"errors"
	"fmt"
)

var ErrArgumentNil = errors.New("argument must not be nil value")

func ErrNotFound(object interface{}) error {
	return fmt.Errorf("not found: %v", object)
}
