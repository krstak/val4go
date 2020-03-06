package val4go

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type validation struct {
	name     string
	validate func(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error
}

func valRequired(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	if vf.Kind() == reflect.Ptr && vf.IsNil() {
		return fmt.Errorf("%s is required", sf.Name)
	}

	if vf.String() == "" {
		return fmt.Errorf("%s is required", sf.Name)
	}
	return nil
}

func valNotBlank(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	if strings.TrimSpace(vf.String()) == "" {
		return fmt.Errorf("%s must not be blank", sf.Name)
	}
	return nil
}

func valEmail(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	r := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !r.MatchString(vf.String()) {
		return fmt.Errorf("%s is not valid email", sf.Name)
	}

	return nil
}

func valCrossEqualField(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	if !reflect.DeepEqual(vf.Interface(), reflect.Indirect(v).FieldByName(cf).Interface()) {
		return fmt.Errorf("%s doesn't match %s", sf.Name, cf)
	}
	return nil
}
