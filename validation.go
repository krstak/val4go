package val4go

import (
	"fmt"
	"reflect"
	"strings"
)

type validation struct {
	name     string
	validate func(v reflect.Value, s reflect.StructField) error
}

func valRequired(v reflect.Value, s reflect.StructField) error {
	if v.String() == "" {
		return fmt.Errorf("field %s is required", s.Name)
	}
	return nil
}

func valNotBlank(v reflect.Value, s reflect.StructField) error {
	if strings.TrimSpace(v.String()) == "" {
		return fmt.Errorf("field %s must not be blank", s.Name)
	}
	return nil
}
