package val4go_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/krstak/testify"
	"github.com/krstak/val4go"
)

func TestValidateRequired(t *testing.T) {
	user := struct {
		FirstName string            `my_schema:"required"`
		LastName  *string           `my_schema:"required"`
		Age       *int              `my_schema:"required"`
		Books     []string          `my_schema:"required"`
		Numbers   []string          `my_schema:"required"`
		Books2    map[string]string `my_schema:"required"`
		Numbers2  map[string]string `my_schema:"required"`
	}{Age: intPtr(12), Numbers: []string{}, Numbers2: make(map[string]string)}

	v := val4go.New()
	v.RegisterSchema("my_schema")

	errs := v.Validate("my_schema", user)
	testify.Equal(t)([]error{errors.New("LastName is required"), errors.New("Books is required"), errors.New("Books2 is required")}, errs)
}

func TestValidateNotEmpty(t *testing.T) {
	user := struct {
		FirstName string `my_schema:"notempty"`
		LastName  string `my_schema:"notempty"`
	}{
		FirstName: "John",
	}

	v := val4go.New()
	v.RegisterSchema("my_schema")

	errs := v.Validate("my_schema", user)
	testify.Equal(t)([]error{errors.New("LastName must not be empty")}, errs)
}

func TestValidateEmail(t *testing.T) {
	type user struct {
		Email string `my_schema:"email"`
	}

	tests := []struct {
		u    user
		errs []error
	}{
		{u: user{}, errs: []error{errors.New("Email is not valid email")}},
		{u: user{Email: "testgmail.com"}, errs: []error{errors.New("Email is not valid email")}},
		{u: user{Email: "test@gmail."}, errs: []error{errors.New("Email is not valid email")}},
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

func TestValidateMinimum(t *testing.T) {
	type user struct {
		FirstName string `my_schema:"min=4"`
		Age       int    `my_schema:"min=18"`
	}

	tests := []struct {
		u    user
		errs []error
	}{
		{u: user{}, errs: []error{errors.New("FirstName must be at least 4 characters long"), errors.New("Age must be minimum 18")}},
		{u: user{FirstName: "Per"}, errs: []error{errors.New("FirstName must be at least 4 characters long"), errors.New("Age must be minimum 18")}},
		{u: user{FirstName: "John"}, errs: []error{errors.New("Age must be minimum 18")}},
		{u: user{FirstName: "John", Age: 17}, errs: []error{errors.New("Age must be minimum 18")}},
		{u: user{FirstName: "John", Age: 18}, errs: []error(nil)},
		{u: user{FirstName: "Per", Age: 21}, errs: []error{errors.New("FirstName must be at least 4 characters long")}},
	}

	v := val4go.New()
	v.RegisterSchema("my_schema")

	for _, ts := range tests {
		errs := v.Validate("my_schema", ts.u)
		testify.Equal(t)(ts.errs, errs)
	}
}

func TestValidateMaximum(t *testing.T) {
	type user struct {
		FirstName string `my_schema:"max=4"`
		Age       int    `my_schema:"max=18"`
	}

	tests := []struct {
		u    user
		errs []error
	}{
		{u: user{}, errs: []error(nil)},
		{u: user{FirstName: "Johnn", Age: 17}, errs: []error{errors.New("FirstName must be at most 4 characters long")}},
		{u: user{FirstName: "Johnn", Age: 18}, errs: []error{errors.New("FirstName must be at most 4 characters long")}},
		{u: user{FirstName: "Johnn", Age: 20}, errs: []error{errors.New("FirstName must be at most 4 characters long"), errors.New("Age must be maximum 18")}},
		{u: user{FirstName: "Pern", Age: 21}, errs: []error{errors.New("Age must be maximum 18")}},
	}

	v := val4go.New()
	v.RegisterSchema("my_schema")

	for _, ts := range tests {
		errs := v.Validate("my_schema", ts.u)
		testify.Equal(t)(ts.errs, errs)
	}
}

func TestValidateDate(t *testing.T) {
	type user struct {
		Date  string  `my_schema:"date"`
		Date2 *string `my_schema:"date"`
	}

	date2 := "2021-05-01"

	tests := []struct {
		u    user
		errs []error
	}{
		{u: user{}, errs: []error(nil)},
		{u: user{Date: "2021-13-01"}, errs: []error{errors.New("Date is not valid date")}},
		{u: user{Date: "2021-11-32"}, errs: []error{errors.New("Date is not valid date")}},
		{u: user{Date: "02/12/2020"}, errs: []error{errors.New("Date is not valid date")}},
		{u: user{Date: "2021-12-01", Date2: &date2}, errs: []error(nil)},
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
		{u: user{Password: "123", ConfirmationPassword: "124"}, errs: []error{errors.New("Password doesn't match ConfirmationPassword")}},
		{u: user{Password: " ", ConfirmationPassword: ""}, errs: []error{errors.New("Password doesn't match ConfirmationPassword")}},
		{u: user{Password: "123", ConfirmationPassword: "12"}, errs: []error{errors.New("Password doesn't match ConfirmationPassword")}},
		{u: user{Age: 12, AgeC: 12}, errs: []error(nil)},
		{u: user{Age: 12, AgeC: 13}, errs: []error{errors.New("Age doesn't match AgeC")}},
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
		Age int `my_schema:"min=18,max=19"`
	}

	tests := []struct {
		u    user
		errs []error
	}{
		{u: user{Age: 17}, errs: []error{errors.New("Age must be minimum 18")}},
		{u: user{Age: 18}, errs: []error(nil)},
		{u: user{Age: 19}, errs: []error(nil)},
		{u: user{Age: 20}, errs: []error{errors.New("Age must be maximum 19")}},
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
		FirstName string `my_schema:"min=1" another_schema:"min=3"`
	}

	u := user{FirstName: "M"}
	v := val4go.New()
	v.RegisterSchema("my_schema")
	v.RegisterSchema("another_schema")

	errs := v.Validate("my_schema", u)
	testify.Equal(t)([]error(nil), errs)

	errs = v.Validate("another_schema", u)
	testify.Equal(t)([]error{errors.New("FirstName must be at least 3 characters long")}, errs)
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

func intPtr(x int) *int {
	return &x
}
