package env

// TODO: make a package out of this idea

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func Parse(structPtr interface{}) error {
	v := reflect.ValueOf(structPtr)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return fmt.Errorf("cannot parse non-pointer or nil value from env")
	}

	if err := parseStructFields(v.Elem()); err != nil {
		return fmt.Errorf("unable to parse %s from env due to the following error(s):\n\n%w", v.Elem().Type().Name(), err)
	}

	return nil
}

func MustParse(structPtr interface{}) {
	if err := Parse(structPtr); err != nil {
		panic(err)
	}
}

func parseStructFields(v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return errors.New("only works with structs")
	}

	errs := make([]error, 0, v.NumField())

	for n := 0; n < v.NumField(); n++ {
		ft := v.Type().Field(n)
		fv := v.Field(n)

		key, ok := ft.Tag.Lookup("env")
		if !ok {
			errs = append(errs, fmt.Errorf("tag `env` not found for field %s", ft.Name))
			continue
		}

		if fv.CanSet() {
			val := os.Getenv(strings.ToUpper(key))

			switch ft.Type.Kind() {
			case reflect.String:
				fv.SetString(val)
			case reflect.Bool:
				b, err := strconv.ParseBool(val)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to parse bool for field %s\n\t%w\n", ft.Name, err))
					continue
				}
				fv.SetBool(b)
			case reflect.Int:
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to parse int for field %s\n\t%w\n", ft.Name, err))
					continue
				}
				fv.SetInt(i)
			case reflect.Int64:
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to parse int64 for field %s\n\t%w\n", ft.Name, err))
					continue
				}
				fv.SetInt(i)
			case reflect.Int32:
				i, err := strconv.ParseInt(val, 10, 32)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to parse int32 for field %s\n\t%w\n", ft.Name, err))
					continue
				}
				fv.SetInt(i)
			case reflect.Int16:
				i, err := strconv.ParseInt(val, 10, 16)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to parse int16 for field %s\n\t%w\n", ft.Name, err))
					continue
				}
				fv.SetInt(i)
			case reflect.Int8:
				i, err := strconv.ParseInt(val, 10, 8)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to parse int8 for field %s\n\t%w\n", ft.Name, err))
					continue
				}
				fv.SetInt(i)
			case reflect.Uint:
				i, err := strconv.ParseUint(val, 10, 64)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to parse uint for field %s\n\t%w\n", ft.Name, err))
					continue
				}
				fv.SetUint(i)
			case reflect.Uint64:
				i, err := strconv.ParseUint(val, 10, 64)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to parse uint64 for field %s\n\t%w\n", ft.Name, err))
					continue
				}
				fv.SetUint(i)
			case reflect.Uint32:
				i, err := strconv.ParseUint(val, 10, 32)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to parse uint32 for field %s\n\t%w\n", ft.Name, err))
					continue
				}
				fv.SetUint(i)
			case reflect.Uint16:
				i, err := strconv.ParseUint(val, 10, 16)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to parse uint16 for field %s\n\t%w\n", ft.Name, err))
					continue
				}
				fv.SetUint(i)
			case reflect.Uint8:
				i, err := strconv.ParseUint(val, 10, 8)
				if err != nil {
					errs = append(errs, fmt.Errorf("unable to parse uint8 for field %s\n\t%w\n", ft.Name, err))
					continue
				}
				fv.SetUint(i)
			default:
				panic("unimplemented")
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	} else {
		return nil
	}
}
