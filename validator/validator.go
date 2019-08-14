package validator

import (
	"reflect"
	"strings"

	validator_v9 "gopkg.in/go-playground/validator.v9"
)

var validate *validator_v9.Validate

func init() {
	validate = validator_v9.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})
}

// New returns a new instance of Validate with sane defaults.
func New() *validator_v9.Validate {
	return validator_v9.New()
}

// Instance returns the Validate instance.
func Instance() *validator_v9.Validate {
	return validate
}

// Validate is a shortcut for Instance().Struct()
func Validate(s interface{}) error {
	return validate.Struct(s)
}
