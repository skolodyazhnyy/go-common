package logger

import (
	"reflect"
	"testing"
)

func TestMerge(t *testing.T) {
	want := map[string]interface{}{
		"foo": "bar",
		"baz": "qux",
		"bar": "bar",
	}

	got := merge(
		map[string]interface{}{"foo": "qux", "baz": "qux"},
		map[string]interface{}{"foo": "bar", "bar": "bar"},
		nil,
	)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Merged value does not match, want: %#v, got: %#v", want, got)
	}
}
