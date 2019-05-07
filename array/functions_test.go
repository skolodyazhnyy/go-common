package array

import (
	"reflect"
	"testing"
)

func TestDifferenceCaseOne(t *testing.T) {

	a := []string{"A", "B", "C"}
	b := []string{"A"}
	expected := []string{"B", "C"}

	res := Difference(a, b)

	if !reflect.DeepEqual(expected, res) {
		t.Errorf("expected %v is not what has been returned %v", expected, res)
	}
}

func TestDifferenceCaseTwo(t *testing.T) {

	a := []string{"A", "B", "C"}
	b := []string{"D"}
	expected := []string{"A", "B", "C"}

	res := Difference(a, b)

	if !reflect.DeepEqual(expected, res) {
		t.Errorf("expected %v is not what has been returned %v", expected, res)
	}
}
