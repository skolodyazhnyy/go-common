package validator

import "reflect"

type validator func([]string, reflect.Value, reflect.StructTag) error

// list of globally defined validators
var validators = []validator{
	required,
}

// required validator checks if field is not empty, it's activated using `required:"true"` structure tag
func required(path []string, val reflect.Value, tag reflect.StructTag) error {
	if val, ok := tag.Lookup("required"); !ok || val != "true" {
		return nil
	}

	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return Error{Path: path, Message: "required field is empty"}
		}
	case reflect.String:
		if val.String() == "" {
			return Error{Path: path, Message: "required field is empty"}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if val.Int() == 0 {
			return Error{Path: path, Message: "required field is empty"}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if val.Uint() == 0 {
			return Error{Path: path, Message: "required field is empty"}
		}
	case reflect.Float32, reflect.Float64:
		if val.Float() == 0 {
			return Error{Path: path, Message: "required field is empty"}
		}
	}

	return nil
}
