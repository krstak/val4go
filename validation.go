package val4go

import (
	"fmt"
	"reflect"
	"regexp"
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

func valEmail(v reflect.Value, s reflect.StructField) error {
	r := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !r.MatchString(v.String()) {
		return fmt.Errorf("field %s is not valid email", s.Name)
	}

	return nil
}
