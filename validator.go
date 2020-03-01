package val4go

import (
	"reflect"
	"strings"
	"sync"
)

type Validator struct {
	schemas     []string
	validations []validation
	sm          sync.RWMutex
	vm          sync.RWMutex
}

func New() *Validator {
	return &Validator{
		schemas: []string{},
		validations: []validation{
			{"required", valRequired},
			{"notblank", valNotBlank},
		},
	}
}

func (v *Validator) RegisterSchema(schema string) {
	v.sm.Lock()
	defer v.sm.Unlock()

	v.schemas = append(v.schemas, schema)
}

func (v *Validator) Validate(schema string, s interface{}) []error {
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
			for _, val := range vals {
				if validation.name == strings.TrimSpace(val) {
					err := validation.validate(valueField, typeField)
					if err != nil {
						errs = append(errs, err)
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
