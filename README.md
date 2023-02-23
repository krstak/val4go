# Struct validation library

_ATTENTION: work in progress_

## Purpose

A small validation library that allows to have multiple different validations on the same struct.

## Usage

Example: a struct used for both cases, sign up and sign in:

```go
type UserAuth struct {
	Email                string `signup:"email" signin:"email"`
	Password             string `signup:"notempty,eq=PasswordConfirmation" signin:"notempty"`
	PasswordConfirmation string `signup:"notempty"`
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

#### reguired

Checks to not be nil
Only `pointers`, `slice` and `map` are supported.
```go
Name *string `signup:"reguired"`
```

#### notempty

Checks to be non-zero values. Only `string` and `slice` are supported.

If field is a pointer, it checks first to not be nil
```go
Name string `signup:"notempty"`
```

#### email

Checks to be a valid email.

If field is a pointer, it checks first to not be nil
```go
Email string `signup:"email"`
```

#### min

Checks a minimum length/value. 

Only `string`, `int`, `int8`, `int16`, `int32`, `int64` are supported.

If field is a pointer, it checks first to not be nil
```go
Name string `signup:"min=4"`
```

#### max

Checks a maximum length/value. 

Only `string`, `int`, `int8`, `int16`, `int32`, `int64` are supported.

If field is a pointer, it checks first to not be nil
```go
Name string `signup:"max=4"`
```

#### date

Checks if it's an iso date.

```go
Date string `signup:"date"`
```

### Cross-field validations:

#### eq

Checks an equality with the cross given field value. 

```go
Password             string `signup:"eq=PasswordConfirmation"`
PasswordConfirmation string
```

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