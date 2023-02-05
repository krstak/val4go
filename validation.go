package val4go

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type validation struct {
	name     string
	validate func(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error
}

func valRequired(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	kind := vf.Kind()

	ptr := kind == reflect.Ptr || kind == reflect.Slice || kind == reflect.Map

	if ptr && vf.IsNil() {
		return fmt.Errorf("%s is required", sf.Name)
	}

	return nil
}

func valNotEmpty(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	err := valRequired(vf, sf, v, cf)
	if err != nil {
		return err
	}

	_, valuefield := value(vf)

	switch vf.Kind() {
	case reflect.String:
		if strings.TrimSpace(valuefield.String()) == "" {
			return fmt.Errorf("%s must not be empty", sf.Name)
		}
	}

	return nil
}

func valEmail(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	err := valRequired(vf, sf, v, cf)
	if err != nil {
		return err
	}

	_, valuefield := value(vf)

	r := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if !r.MatchString(valuefield.String()) {
		return fmt.Errorf("%s is not valid email", sf.Name)
	}

	return nil
}

func valMin(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	err := valRequired(vf, sf, v, cf)
	if err != nil {
		return err
	}

	//todo: convert to int64
	min, err := strconv.Atoi(cf)
	if err != nil {
		return err
	}

	kind, valuefield := value(vf)

	switch kind {
	case reflect.String:
		if len(valuefield.String()) < min {
			return fmt.Errorf("%s must be at least %s characters long", sf.Name, cf)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if valuefield.Int() < int64(min) {
			return fmt.Errorf("%s must be minimum %s", sf.Name, cf)
		}
	}

	return nil
}

func valMax(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	err := valRequired(vf, sf, v, cf)
	if err != nil {
		return err
	}

	//todo: convert to int64
	max, err := strconv.Atoi(cf)
	if err != nil {
		return err
	}

	kind, valuefield := value(vf)

	switch kind {
	case reflect.String:
		if len(valuefield.String()) > max {
			return fmt.Errorf("%s must be at most %s characters long", sf.Name, cf)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if valuefield.Int() > int64(max) {
			return fmt.Errorf("%s must be maximum %s", sf.Name, cf)
		}
	}

	return nil
}

func valCrossEqualField(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	err := valRequired(vf, sf, v, cf)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(vf.Interface(), reflect.Indirect(v).FieldByName(cf).Interface()) {
		return fmt.Errorf("%s doesn't match %s", sf.Name, cf)
	}
	return nil
}

func valDate(vf reflect.Value, sf reflect.StructField, v reflect.Value, cf string) error {
	date := vf.String()
	if date == "" {
		return nil
	}

	kind := vf.Kind()
	if kind == reflect.Ptr && vf.IsNil() {
		return nil
	}

	_, valuefield := value(vf)

	if _, err := time.Parse("2006-01-02", valuefield.String()); err != nil {
		return fmt.Errorf("%s is not valid date", sf.Name)
	}

	return nil
}

func value(vf reflect.Value) (reflect.Kind, reflect.Value) {
	kind := vf.Kind()
	valuefield := vf
	if kind == reflect.Ptr {
		kind = vf.Elem().Kind()
		valuefield = vf.Elem()
	}

	return kind, valuefield
}
