package walker

import "reflect"

type callback func(v interface{}, sf reflect.StructField)

// Walk will apply a given "callback" to each element in a given structure recursively
// 	v := &T
// 	f := callback
func Walk(v interface{}, f callback) {
	// Reflect & dereference the receiver structure pointer
	o := reflect.ValueOf(v).Elem()

	// Start from whatever is the 1st field in structure.
	t := o.Type().Field(0)

	// Start traversing...
	recurse(o, t, f)
}

// recurse does all the internal magic of traversing the structure and applying
// the callback identified by "f"
func recurse(obj reflect.Value, sf reflect.StructField, f callback) {
	switch obj.Kind() {
	// Dereference pointer and recurse...
	case reflect.Ptr:
		val := obj.Elem()
		if !val.IsValid() {
			return
		}
		recurse(val, sf, f)

	// Enter the struct
	case reflect.Struct:
		for i := 0; i < obj.NumField(); i++ {
			fld := obj.Type().Field(i)
			recurse(obj.Field(i), fld, f)
		}

	// And everything else will have the callback applied
	default:
		f(&obj, sf)
	}
}
