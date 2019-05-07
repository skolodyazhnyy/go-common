package slice

import (
	"reflect"
	"testing"
)

func TestDifferenceCaseOne(t *testing.T) {

	var difftests = []struct {
		name     string
		in1      []string
		in2      []string
		expected []string
	}{
		{"Empty slices", []string{}, []string{}, []string{}},
		{"Empty slice 1", []string{}, []string{"A", "B", "C"}, []string{}},
		{"Empty slice 2", []string{"A", "B", "C"}, []string{}, []string{"A", "B", "C"}},
		{"Two elements present in slice 1, but not in slice 2", []string{"A", "B", "C"}, []string{"A"}, []string{"B", "C"}},
		{"No element match in slice 2", []string{"A", "B", "C"}, []string{"D"}, []string{"A", "B", "C"}},
		{"Slice 1 and slice 2 are equal", []string{"A", "B", "C"}, []string{"A", "B", "C"}, []string{}},
	}

	for _, tt := range difftests {
		t.Run(tt.name, func(t *testing.T) {
			res := Difference(tt.in1, tt.in2)
			if !reflect.DeepEqual(tt.expected, res) {
				t.Errorf("expected %v is not what has been returned %v", tt.expected, res)
			}
		})
	}
}
