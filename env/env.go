package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/magento-mcom/go-common/walker"
)

// Unmarshal will parse structure s for 'env:"ENV_VAR"' tags and
// will attempt to populate it with values found in the current environment or
// whatever default value is specified.
//
// To unmarshal ENV vars into an interface value, Unmarshal will read the type
// of the specific field in the passed structure to typecast the value read.
func Unmarshal(s interface{}) {
	// v is actually a pointer to the reflect.Value
	walker.Walk(s, func(v interface{}, sf reflect.StructField) {
		tag := sf.Tag.Get("env")
		if tag == "" {
			return
		}

		if value, ok := os.LookupEnv(tag); ok {
			kind := sf.Type.Kind()
			val := v.(*reflect.Value)

			switch kind {
			case reflect.String:
				val.SetString(value)
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				n, err := strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
				val.SetInt(int64(n))
			case reflect.Bool:
				n, err := strconv.Atoi(value)
				if err != nil {
					panic(err)
				}
				val.SetBool(int64(n) != 0)
			case reflect.Float32, reflect.Float64:
				n, err := strconv.ParseFloat(value, 64)
				if err != nil {
					panic(err)
				}
				val.SetFloat(float64(n))
			default:
				panic(fmt.Sprintf("Unsupported type %s for field %s", kind, sf.Name))
			}
		}
	})
}
