package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
)

type RegexpValidator struct {
	BaseValidator
}

func NewRegexpValidator() RegexpValidator {
	types := make(TypeKindDict, 1)
	types[reflect.String] = true
	return RegexpValidator{
		BaseValidator{
			name:         regexpValidator,
			allowedTypes: types,
		},
	}
}

func (rxpv RegexpValidator) Validate(
	fieldName string,
	reflectVal reflect.Value,
	cond validationCondition,
) (ValidationError, bool) {
	valErr := ValidationError{Field: fieldName}

	compiled := regexp.MustCompile(cond.rule)
	ok := compiled.MatchString(reflectVal.String())
	if !ok {
		valErr.Err = fmt.Errorf(ErrValidationStringRegexp.Error(), fieldName, cond.rule)
		return valErr, false
	}

	return valErr, true
}
