package val4go

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
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

func valMin(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	min, err := strconv.Atoi(cf)
	if err != nil {
		return err
	}
	switch vf.Kind() {
	case reflect.String:
		if len(vf.String()) < min {
			return fmt.Errorf("%s must be at least %s characters long", sf.Name, cf)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if vf.Int() < int64(min) {
			return fmt.Errorf("%s must be minimum %s", sf.Name, cf)
		}
	}

	return nil
}

func valCrossEqualField(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	if !reflect.DeepEqual(vf.Interface(), reflect.Indirect(v).FieldByName(cf).Interface()) {
		return fmt.Errorf("%s doesn't match %s", sf.Name, cf)
	}
	return nil
}
