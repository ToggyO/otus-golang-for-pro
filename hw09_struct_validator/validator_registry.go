package hw09structvalidator

import "reflect"

type TypeKindDict map[reflect.Kind]bool

type validationCondition struct {
	name string
	rule string
}

type IValidator interface {
	Validate(fieldName string, reflectVal reflect.Value, cond validationCondition) (ValidationError, bool)
	Tag() string
	CanValidate(t reflect.Type) bool
}

type ValidatorRegistry map[string]IValidator

func NewValidatorRegistry(validators []IValidator) ValidatorRegistry {
	reg := make(ValidatorRegistry, len(validators))
	for _, v := range validators {
		reg[v.Tag()] = v
	}
	return reg
}
