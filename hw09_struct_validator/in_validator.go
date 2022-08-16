package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type InValidator struct {
	BaseValidator
}

func NewInValidator() InValidator {
	types := make(TypeKindDict, 6)
	types[reflect.String] = true
	types[reflect.Int64] = true
	types[reflect.Int32] = true
	types[reflect.Int16] = true
	types[reflect.Int8] = true
	types[reflect.Int] = true

	return InValidator{
		BaseValidator{
			name:         inValidator,
			allowedTypes: types,
		},
	}
}

func (mv InValidator) Validate(
	fieldName string,
	reflectVal reflect.Value,
	cond validationCondition,
) (ValidationError, bool) {
	valErr := ValidationError{Field: fieldName}

	var strVal string
	switch reflectVal.Kind() { //nolint:exhaustive
	case reflect.String:
		strVal = reflectVal.String()
	case reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
		strVal = fmt.Sprint(reflectVal.Int())
	default:
		valErr.Err = fmt.Errorf(ErrUnsupportedType.Error(), reflectVal.Type())
		return valErr, false
	}

	if len(strVal) == 0 || !strings.Contains(cond.rule, strVal) {
		valErr.Err = fmt.Errorf(ErrValidationIn.Error(), fieldName, cond.rule)
		return valErr, false
	}

	return valErr, true
}
