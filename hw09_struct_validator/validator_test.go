package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`

		Application App             `validate:"nested"`
		meta        json.RawMessage //nolint:structcheck,unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
		IsPublic  byte `validate:"max:256"`
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
		Data []int  `validate:"min:10|max:15"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "common struct validation with nested struct",
			in: User{
				ID:     "23817387219379",
				Name:   "",
				Age:    17,
				Email:  "kek@mail.ru",
				Phones: []string{"89134565432"},
				Application: App{
					Version: "1.2.0-rc1",
				},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: fmt.Errorf(ErrValidationStringLength.Error(), "ID", 36)},
				ValidationError{Field: "Age", Err: fmt.Errorf(ErrValidationMin.Error(), "Age", 18)},
				ValidationError{Field: "Role", Err: fmt.Errorf(ErrValidationIn.Error(), "Role", "admin,stuff")},
				ValidationError{Field: "Version", Err: fmt.Errorf(ErrValidationStringLength.Error(), "Version", 5)},
			},
		},
		{
			name: "slice validation error",
			in: Response{
				Code: 404,
				Body: "",
				Data: []int{5, 20, 11},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Data", Err: fmt.Errorf(ErrValidationMin.Error(), "Data", 10)},
				ValidationError{Field: "Data", Err: fmt.Errorf(ErrValidationMax.Error(), "Data", 15)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			var valErrs ValidationErrors
			require.ErrorAs(t, err, &valErrs)

			for _, e := range strings.Split(tt.expectedErr.Error(), "\n") {
				require.ErrorContains(t, err, e)
			}

			_ = tt
		})
	}
}

func TestUnsupportedType(t *testing.T) {
	var isPublic byte = 45
	token := Token{
		Header:    []byte("Header"),
		Payload:   []byte("Payload"),
		Signature: []byte("Signature"),
		IsPublic:  isPublic,
	}

	t.Run("unsupported data type", func(t *testing.T) {
		err := Validate(token)
		require.EqualError(t, err, fmt.Errorf(ErrUnsupportedType.Error(), reflect.TypeOf(isPublic)).Error())
	})
}
