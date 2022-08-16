package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
)

type MinValidator struct {
	BaseValidator
}

func NewMinValidator() MinValidator {
	types := make(TypeKindDict, 5)
	types[reflect.Int64] = true
	types[reflect.Int32] = true
	types[reflect.Int16] = true
	types[reflect.Int8] = true
	types[reflect.Int] = true

	return MinValidator{
		BaseValidator{
			name:         minValidator,
			allowedTypes: types,
		},
	}
}

func (mv MinValidator) Validate(
	fieldName string,
	reflectVal reflect.Value,
	cond validationCondition,
) (ValidationError, bool) {
	valErr := ValidationError{Field: fieldName}

	num, err := strconv.Atoi(cond.rule)
	if reflectVal.Int() < int64(num) || err != nil {
		valErr.Err = fmt.Errorf(ErrValidationMin.Error(), fieldName, num)
		return valErr, false
	}

	return valErr, true
}
