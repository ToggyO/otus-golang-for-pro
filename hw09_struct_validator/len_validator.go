package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
)

type LenValidator struct {
	BaseValidator
}

func NewLenValidator() LenValidator {
	types := make(TypeKindDict, 1)
	types[reflect.String] = true

	return LenValidator{
		BaseValidator{
			name:         lenValidator,
			allowedTypes: types,
		},
	}
}

func (lv LenValidator) Validate(
	fieldName string,
	reflectVal reflect.Value,
	cond validationCondition,
) (ValidationError, bool) {
	valErr := ValidationError{Field: fieldName}

	str := reflectVal.String()
	expectedLength, err := strconv.Atoi(cond.rule)
	if len(str) != expectedLength || err != nil {
		valErr.Err = fmt.Errorf(ErrValidationStringLength.Error(), fieldName, expectedLength)
		return valErr, false
	}

	return valErr, true
}
