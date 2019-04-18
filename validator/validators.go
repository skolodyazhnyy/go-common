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
	}

	return nil
}
