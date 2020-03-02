package val4go

import (
	"reflect"
	"strings"
	"sync"
)

type V interface {
	RegisterSchema(string)
	RegisterValidation(string, func(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error)
	Validate(string, interface{}) []error
}

type validator struct {
	schemas     []string
	validations []validation
	sm          sync.RWMutex
	vm          sync.RWMutex
}

func New() V {
	v := &validator{
		schemas:     []string{},
		validations: make([]validation, 0, 4),
	}

	v.RegisterValidation("required", valRequired)
	v.RegisterValidation("notblank", valNotBlank)
	v.RegisterValidation("email", valEmail)
	v.RegisterValidation("eq", valCrossEqualField)

	return v
}

func (v *validator) RegisterSchema(schema string) {
	v.vm.Lock()
	defer v.vm.Unlock()

	v.schemas = append(v.schemas, schema)
}

func (v *validator) RegisterValidation(name string, fn func(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error) {
	v.sm.Lock()
	defer v.sm.Unlock()

	v.validations = append(v.validations, validation{name: name, validate: fn})
}

func (v *validator) Validate(schema string, s interface{}) []error {
	v.sm.RLock()
	schemas := v.schemas
	v.sm.RUnlock()

	v.vm.RLock()
	validations := v.validations
	v.vm.RUnlock()

	return validate(validations, schemas, schema, s)
}

func validate(validations []validation, schemas []string, schema string, s interface{}) []error {
	errs := []error(nil)
	val := reflect.ValueOf(s)

	if !contains(schemas, schema) {
		return errs
	}

	if val.Kind() != reflect.Struct {
		return errs
	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		name := typeField.Tag.Get(schema)
		if name == "" {
			continue
		}

		for _, validation := range validations {
			vals := strings.Split(name, ",")
			for _, vld := range vals {
				if v, f, ok := cross(vld); ok {
					if validation.name == v {
						err := validation.validate(valueField, typeField, val, f)
						if err != nil {
							errs = append(errs, err)
						}
					}
				} else {
					if validation.name == strings.TrimSpace(vld) {
						err := validation.validate(valueField, typeField, val, "")
						if err != nil {
							errs = append(errs, err)
						}
					}
				}
			}
		}
	}

	return errs
}

func contains(vals []string, val string) bool {
	for _, s := range vals {
		if s == val {
			return true
		}
	}
	return false
}

func cross(name string) (string, string, bool) {
	parts := strings.Split(name, "=")
	if len(parts) == 2 {
		return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), true
	}
	return "", "", false
}
