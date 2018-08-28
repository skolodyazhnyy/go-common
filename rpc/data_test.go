package rpc

import (
	"reflect"
	"testing"
)

type testData struct {
	id     string
	params map[string]string
}

func TestBinding(t *testing.T) {
	givenID := "bazzinga"
	givenParams := map[string]string{"f": "foo", "q": "qux"}

	b := NewStructData(testData{
		id:     givenID,
		params: givenParams,
	})

	payload := &testData{params: map[string]string{}}

	err := b.Bind(payload)

	if err != nil {
		t.Errorf("Binding failed: %v", err)
	}

	if payload.id != givenID {
		t.Errorf("Binding failed, ID does not match: %v instead of %v", payload.id, givenID)
	}

	if !reflect.DeepEqual(payload.params, givenParams) {
		t.Errorf("Binding failed, Params do not match: %v instead of %v", payload.params, givenParams)
	}
}

func TestBindingConvertiblePayload(t *testing.T) {
	givenPayload := int32(19)

	b := NewStructData(givenPayload)

	payload := int64(0)
	err := b.Bind(&payload)

	if err != nil {
		t.Errorf("Binding failed: %v", err)
	}

	if int32(payload) != givenPayload {
		t.Errorf("Binding failed, value does not match: %v instead of %v", payload, givenPayload)
	}
}

func TestBindingTypeMismatch(t *testing.T) {
	b := NewStructData("boom")

	payload := &testData{params: map[string]string{}}

	err := b.Bind(payload)

	if err == nil {
		t.Error("Binding suppose to fail but it didn't")
	}
}
