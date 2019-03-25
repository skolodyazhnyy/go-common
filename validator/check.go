package validator

import (
	"fmt"
	"reflect"
)

// Check if given data structure is valid
func Check(in interface{}) error {
	return traverse(nil, reflect.ValueOf(in))
}

// traverse given value checking each structure field, element of the slice or map
func traverse(path []string, val reflect.Value) (errs Errors) {
	switch val.Kind() {
	case reflect.Ptr:
		return traverse(path, val.Elem())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			fv := val.Field(i)
			ft := val.Type().Field(i)
			fp := append(path, ft.Name)

			errs = append(errs, validate(fp, fv, ft.Tag)...)

			if fv.IsValid() {
				errs = append(errs, traverse(fp, fv)...)
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			errs = append(errs, traverse(append(path, fmt.Sprint(i)), val.Index(i))...)
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			errs = append(errs, traverse(append(path, key.String()), val.MapIndex(key))...)
		}
	}

	return
}

// validate field of the structure against predefined list of validators
func validate(path []string, v reflect.Value, t reflect.StructTag) (errs Errors) {
	for _, validate := range validators {
		if err := validate(path, v, t); err != nil {
			errs = append(errs, err)
		}
	}

	return
}
