package walker

import (
	"fmt"
	"reflect"
	"strings"
)

type structure struct {
	A      string
	B      int
	C      bool
	D      float64
	Level1 level1
}

type level1 struct {
	H     string
	Inner level2
}

type level2 struct {
	X string
	Z int
}

func ExampleWalk_display() {
	test := structure{
		A: "String A",
		B: 167,
		C: true,
		D: 12.45,
		Level1: level1{
			H: "String H",
			Inner: level2{
				X: "X",
				Z: 1987,
			},
		},
	}

	Walk(&test, func(v interface{}, sf reflect.StructField) {
		fmt.Println(v.(*reflect.Value).String())
	})

	// Output:
	// String A
	// <int Value>
	// <bool Value>
	// <float64 Value>
	// String H
	// X
	// <int Value>
}

func ExampleWalk_modify() {
	test := struct {
		A string
		B string
		C string
	}{
		A: "small a",
		B: "small b",
		C: "small c",
	}

	fmt.Printf("BEFORE: %+v\n", test)

	Walk(&test, func(v interface{}, sf reflect.StructField) {
		val := v.(*reflect.Value)
		upper := strings.ToUpper(val.String())
		val.SetString(upper)
	})

	fmt.Printf("AFTER: %+v\n", test)

	// Output:
	// BEFORE: {A:small a B:small b C:small c}
	// AFTER: {A:SMALL A B:SMALL B C:SMALL C}
}
