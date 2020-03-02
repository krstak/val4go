package val4go_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/krstak/testify"
	"github.com/krstak/val4go"
)

func TestValidateRequired(t *testing.T) {
	type user struct {
		FirstName string `my_schema:"required"`
	}

	tests := []struct {
		u    user
		errs []error
	}{
		{u: user{}, errs: []error{errors.New("field FirstName is required")}},
		{u: user{FirstName: " "}, errs: []error(nil)},
		{u: user{FirstName: "John"}, errs: []error(nil)},
	}

	v := val4go.New()
	v.RegisterSchema("my_schema")

	for _, ts := range tests {
		errs := v.Validate("my_schema", ts.u)
		testify.Equal(t)(ts.errs, errs)
	}
}

func TestValidateRequiredPtr(t *testing.T) {
	type user struct {
		Age *int `my_schema:"required"`
	}

	tests := []struct {
		u    user
		errs []error
	}{
		{u: user{}, errs: []error{errors.New("field Age is required")}},
		{u: user{Age: intP(12)}, errs: []error(nil)},
	}

	v := val4go.New()
	v.RegisterSchema("my_schema")

	for _, ts := range tests {
		errs := v.Validate("my_schema", ts.u)
		testify.Equal(t)(ts.errs, errs)
	}
}

func TestValidateNotBlank(t *testing.T) {
	type user struct {
		FirstName string `my_schema:"notblank"`
	}

	tests := []struct {
		u    user
		errs []error
	}{
		{u: user{}, errs: []error{errors.New("field FirstName must not be blank")}},
		{u: user{FirstName: " "}, errs: []error{errors.New("field FirstName must not be blank")}},
		{u: user{FirstName: "John"}, errs: []error(nil)},
	}

	v := val4go.New()
	v.RegisterSchema("my_schema")

	for _, ts := range tests {
		errs := v.Validate("my_schema", ts.u)
		testify.Equal(t)(ts.errs, errs)
	}
}

func TestValidateEmail(t *testing.T) {
	type user struct {
		Email string `my_schema:"email"`
	}

	tests := []struct {
		u    user
		errs []error
	}{
		{u: user{}, errs: []error{errors.New("field Email is not valid email")}},
		{u: user{Email: "testgmail.com"}, errs: []error{errors.New("field Email is not valid email")}},
		{u: user{Email: "test@gmail."}, errs: []error{errors.New("field Email is not valid email")}},
		{u: user{Email: "test@gmail.com"}, errs: []error(nil)},
		{u: user{Email: "test@gmail"}, errs: []error(nil)},
	}

	v := val4go.New()
	v.RegisterSchema("my_schema")

	for _, ts := range tests {
		errs := v.Validate("my_schema", ts.u)
		testify.Equal(t)(ts.errs, errs)
	}
}

func TestValidateCrossEqfield(t *testing.T) {
	type user struct {
		Password             string `my_schema:"eq=ConfirmationPassword"`
		ConfirmationPassword string
		Age                  int `my_schema:"eq=AgeC"`
		AgeC                 int
	}

	tests := []struct {
		u    user
		errs []error
	}{
		{u: user{}, errs: []error(nil)},
		{u: user{Password: "123", ConfirmationPassword: "123"}, errs: []error(nil)},
		{u: user{Password: " ", ConfirmationPassword: " "}, errs: []error(nil)},
		{u: user{Password: "123", ConfirmationPassword: "124"}, errs: []error{errors.New("field Password doesn't match field ConfirmationPassword")}},
		{u: user{Password: " ", ConfirmationPassword: ""}, errs: []error{errors.New("field Password doesn't match field ConfirmationPassword")}},
		{u: user{Password: "123", ConfirmationPassword: "12"}, errs: []error{errors.New("field Password doesn't match field ConfirmationPassword")}},
		{u: user{Age: 12, AgeC: 12}, errs: []error(nil)},
		{u: user{Age: 12, AgeC: 13}, errs: []error{errors.New("field Age doesn't match field AgeC")}},
	}

	v := val4go.New()
	v.RegisterSchema("my_schema")

	for _, ts := range tests {
		errs := v.Validate("my_schema", ts.u)
		testify.Equal(t)(ts.errs, errs)
	}
}

func TestValidateMultipleValidations(t *testing.T) {
	type user struct {
		FirstName string `my_schema:"required,notblank"`
	}

	tests := []struct {
		u    user
		errs []error
	}{
		{u: user{}, errs: []error{errors.New("field FirstName is required"), errors.New("field FirstName must not be blank")}},
		{u: user{FirstName: " "}, errs: []error{errors.New("field FirstName must not be blank")}},
		{u: user{FirstName: "John"}, errs: []error(nil)},
	}

	v := val4go.New()
	v.RegisterSchema("my_schema")

	for _, ts := range tests {
		errs := v.Validate("my_schema", ts.u)
		testify.Equal(t)(ts.errs, errs)
	}
}

func TestValidateUnregisteredSchema(t *testing.T) {
	type user struct {
		FirstName string `unregistered_schema:"required"`
	}

	u := user{}
	v := val4go.New()

	errs := v.Validate("unregistered_schema", u)
	testify.Equal(t)([]error(nil), errs)
}

func TestValidateDifferentSchema(t *testing.T) {
	type user struct {
		FirstName string `my_schema:"required"`
	}

	u := user{}
	v := val4go.New()
	v.RegisterSchema("my_schema")

	errs := v.Validate("another_schema", u)
	testify.Equal(t)([]error(nil), errs)
}

func TestValidateCustomValidation(t *testing.T) {
	type user struct {
		FirstName string `my_schema:"custom"`
	}

	u := user{FirstName: "John"}
	v := val4go.New()
	v.RegisterSchema("my_schema")

	errs := v.Validate("my_schema", u)
	testify.Equal(t)([]error(nil), errs)

	v.RegisterValidation("custom", func(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
		if vf.String() == "John" {
			return errors.New("john is not valid")
		}
		return nil
	})

	errs = v.Validate("my_schema", u)
	testify.Equal(t)([]error{errors.New("john is not valid")}, errs)
}

func intP(x int) *int {
	return &x
}
