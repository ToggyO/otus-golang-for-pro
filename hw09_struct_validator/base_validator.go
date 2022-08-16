package hw09structvalidator

import "reflect"

type BaseValidator struct {
	name         string
	allowedTypes TypeKindDict
}

func (v BaseValidator) Validate(fieldName string,
	reflectVal reflect.Value,
	cond validationCondition,
) (ValidationError, bool) {
	return ValidationError{}, false
}

func (v BaseValidator) Tag() string {
	return v.name
}

func (v BaseValidator) CanValidate(t reflect.Type) bool {
	return v.allowedTypes[t.Kind()]
}
