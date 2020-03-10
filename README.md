# Struct validation library

_ATTENTION: work in progress_

## Purpose

A small validation library that allows to have multiple different validations on the same struct.

## Usage

Example: a struct used for both cases, sign up and sign in:

```go
type UserAuth struct {
	Email                string `signup:"email" signin:"email"`
	Password             string `signup:"notblank,eq=PasswordConfirmation" signin:"notblank"`
	PasswordConfirmation string `signup:"notblank"`
}

usr := UserAuth{...}
```

Init validator and register schemas:

```go
v := val4go.New()
v.RegisterSchema("signup")
v.RegisterSchema("signin")
```

Validate based on use case:

```go
errs := v.Validate("signup", usr)
```

or

```go
errs := v.Validate("signin", usr)
```

## Validations

### Field validations:
- reguired
- notblank
- email
- min

### Cross-field validations:
- eq=_another-field-name_

## Custom validations:

```go
v := val4go.New()

v.RegisterValidation("my-validation", func(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	// vf -> reflect.Value of a field being validated
	// sf -> reflect.StructField of a field being validated
	// v  -> reflect.Value of a struct being validated
	// cf -> string of a cross field name

	// validate
	return err
})

type UserAuth struct {
	Email                string `signup:"my-validation"`
}

v.RegisterSchema("signup")

```