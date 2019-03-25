package validator

import (
	"reflect"
	"testing"
)

type flat struct {
	Required string `required:"true"`
	Optional string `required:"false"`
}

type nested struct {
	Required *nested `required:"true"`
	Optional *nested `required:"false"`
}

type mapped struct {
	Required map[string]*nested `required:"true"`
	Optional map[string]*nested `required:"false"`
}

type sliced struct {
	Required []*nested `required:"true"`
	Optional []*nested `required:"false"`
}

func TestCheck(t *testing.T) {
	tests := []struct {
		Name   string
		Input  interface{}
		Output Errors
	}{
		{
			Name:   "nil value",
			Input:  nil,
			Output: nil,
		},
		{
			Name:  "flat structure",
			Input: flat{},
			Output: Errors{
				Error{Path: []string{"Required"}, Message: "required field is empty"},
			},
		},
		{
			Name:  "nested structure",
			Input: nested{},
			Output: Errors{
				Error{Path: []string{"Required"}, Message: "required field is empty"},
			},
		},
		{
			Name: "deep nested structure",
			Input: nested{
				Required: &nested{},
			},
			Output: Errors{
				Error{Path: []string{"Required", "Required"}, Message: "required field is empty"},
			},
		},
		{
			Name:   "map structure",
			Input:  mapped{},
			Output: nil,
		},
		{
			Name: "deep map structure",
			Input: mapped{
				Required: map[string]*nested{"KEY": {}},
			},
			Output: Errors{
				Error{Path: []string{"Required", "KEY", "Required"}, Message: "required field is empty"},
			},
		},
		{
			Name:   "slice structure",
			Input:  sliced{},
			Output: nil,
		},
		{
			Name: "deep slice structure",
			Input: sliced{
				Required: []*nested{{}},
			},
			Output: Errors{
				Error{Path: []string{"Required", "0", "Required"}, Message: "required field is empty"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := Check(test.Input)

			errs, ok := err.(Errors)
			if !ok {
				t.Fatalf("Returned value is not of type validator.Errors, got %T instead", err)
			}

			if !reflect.DeepEqual(errs, test.Output) {
				t.Errorf("Errors do not match: want %#v, got %#v", test.Output, errs)
			}
		})
	}
}
