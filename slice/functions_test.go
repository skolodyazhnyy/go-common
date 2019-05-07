package slice

import (
	"reflect"
	"testing"
)

var difftests = []struct {
	in1      []string
	in2      []string
	expected []string
}{
	{[]string{}, []string{}, []string{}},
	{[]string{}, []string{"A", "B", "C"}, []string{}},
	{[]string{"A", "B", "C"}, []string{"A"}, []string{"B", "C"}},
	{[]string{"A", "B", "C"}, []string{"D"}, []string{"A", "B", "C"}},
	{[]string{"A", "B", "C"}, []string{}, []string{"A", "B", "C"}},
	{[]string{"A", "B", "C"}, []string{"A", "B", "C"}, []string{}},
}

func TestDifferenceCaseOne(t *testing.T) {

	for _, tt := range difftests {
		res := Difference(tt.in1, tt.in2)
		if !reflect.DeepEqual(tt.expected, res) {
			t.Errorf("expected %v is not what has been returned %v", tt.expected, res)
		}
	}
}
