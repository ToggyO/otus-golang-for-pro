package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func Validate(v interface{}) error {
	in := reflect.ValueOf(v)
	if in.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var validationErrors ValidationErrors
	refType := reflect.TypeOf(v)
	numFields := refType.NumField()

	for i := 0; i < numFields; i++ {
		field := refType.Field(i)
		if !field.IsExported() {
			continue
		}

		if tag, ok := field.Tag.Lookup(validateTag); ok {
			var err error

			switch field.Type.Kind() { //nolint:exhaustive
			case reflect.String, reflect.Int64, reflect.Int32, reflect.Int16, reflect.Int8, reflect.Int:
				err = validatePrimitive(tag, field.Name, in.Field(i))
			case reflect.Slice, reflect.Array:
				err = validateSlice(tag, field.Name, in.Field(i))
			case reflect.Struct:
				if !strings.Contains(tag, nested) {
					continue
				}
				err = Validate(in.Field(i).Interface())
			default:
				return fmt.Errorf(ErrUnsupportedType.Error(), field.Type)
			}

			var valErrs ValidationErrors
			if errors.As(err, &valErrs) {
				if len(valErrs) > 0 {
					validationErrors = append(validationErrors, valErrs...)
				}
				continue
			}
			return err
		}
	}

	return validationErrors
}

func validatePrimitive(tagValue, fieldName string, val reflect.Value) error {
	var validationErrors ValidationErrors

	conditions := parseTag(tagValue)
	for _, cond := range conditions {
		validator, err := getValidator(cond.name)
		if err != nil {
			return err
		}

		if !validator.CanValidate(val.Type()) {
			continue
		}

		valErr, isValid := validator.Validate(fieldName, val, cond)
		if !isValid {
			validationErrors = append(validationErrors, valErr)
		}
	}

	return validationErrors
}

func validateSlice(tagValue, fieldName string, val reflect.Value) error {
	var validationErrors ValidationErrors

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		err := validatePrimitive(tagValue, fieldName, item)
		var valErrs ValidationErrors
		if errors.As(err, &valErrs) {
			if len(valErrs) > 0 {
				validationErrors = append(validationErrors, valErrs...)
			}
			continue
		}
		return err
	}

	return validationErrors
}

var registry = NewValidatorRegistry([]IValidator{
	NewLenValidator(),
	NewRegexpValidator(),
	NewMinValidator(),
	NewMaxValidator(),
	NewInValidator(),
})

func getValidator(name string) (IValidator, error) {
	validator, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf(ErrUnsupportedValidatorType.Error(), name)
	}

	return validator, nil
}
