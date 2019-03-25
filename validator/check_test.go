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
		Output []error
	}{
		{
			Name:   "nil value",
			Input:  nil,
			Output: nil,
		},
		{
			Name:  "flat structure",
			Input: flat{},
			Output: []error{
				Error{Path: []string{"Required"}, Message: "required field is empty"},
			},
		},
		{
			Name:  "nested structure",
			Input: nested{},
			Output: []error{
				Error{Path: []string{"Required"}, Message: "required field is empty"},
			},
		},
		{
			Name: "deep nested structure",
			Input: nested{
				Required: &nested{},
			},
			Output: []error{
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
			Output: []error{
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
			Output: []error{
				Error{Path: []string{"Required", "0", "Required"}, Message: "required field is empty"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			errs := Check(test.Input)

			if !reflect.DeepEqual(errs, test.Output) {
				t.Errorf("Errors do not match: want %#v, got %#v", test.Output, errs)
			}
		})
	}
}
