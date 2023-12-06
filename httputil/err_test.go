package httputil_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"

	"github.com/avakarev/go-util/httputil"
	"github.com/avakarev/go-util/testutil"
)

func TestNewValidationErr(t *testing.T) {
	type model struct {
		Name string `json:"name" validate:"required"`
		IPV4 string `json:"ipv4" validate:"required,ipv4"`
		Some string `validate:"required"`
	}

	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	var ve validator.ValidationErrors
	testutil.Diff(true, errors.As(validate.Struct(&model{}), &ve), t)

	resp := httputil.NewValidationErr(validator.ValidationErrors(ve))
	testutil.Diff(&httputil.ErrResponse{
		Error: httputil.Err{
			Code: 400,
			Msg:  "validation error",
			Items: []httputil.ValidationErr{
				{Subject: "name", Msg: "required"},
				{Subject: "ipv4", Msg: "required"},
				{Subject: "some", Msg: "required"},
			},
		},
	}, resp, t)
}
