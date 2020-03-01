package val4go_test

import (
	"errors"
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
