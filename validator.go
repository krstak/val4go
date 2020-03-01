package val4go

import (
	"fmt"
	"reflect"
	"sync"
)

type Validator struct {
	schemas []string
	mux     sync.RWMutex
}

func New() *Validator {
	return &Validator{schemas: []string{"val4go"}}
}

func (v *Validator) RegisterSchema(schema string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.schemas = append(v.schemas, schema)
}

func (v *Validator) Validate(schema string, s interface{}) []error {
	errs := []error(nil)

	v.mux.RLock()
	schemas := v.schemas
	v.mux.RUnlock()

	if !contains(schemas, schema) {
		return errs
	}

	val := reflect.ValueOf(s)
	if val.Kind() != reflect.Struct {
		return errs
	}

	for i := 0; i < val.NumField(); i++ {
		fieldInfo := val.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag             // my_schema:"required"

		// fmt.Println(tag)
		// fmt.Println(fieldInfo.Name) // Name

		// a reflect.StructTag
		name := tag.Get(schema) // required
		// fmt.Println(name)
		if name == "" {
			continue
		}

		//todo: check if is a string
		if name == "required" {
			s := val.Field(i).String()
			if s == "" {
				errs = append(errs, fmt.Errorf("field %s is required", fieldInfo.Name))
			}
		}

		// fields[name] = v.Field(i)
	}
	// t.Kind()
	// e := t.Elem()
	// fmt.Println("------------")
	// fmt.Println(elem)
	// // fmt.Println(va.Field(0))
	// fmt.Println("------------")

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
