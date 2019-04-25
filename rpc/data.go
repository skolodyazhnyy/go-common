package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/magento-mcom/go-common/rpc/errors"
	"reflect"
)

// Data is container for a value which is not yet bound to a structure.
// It is used to represent request parameters and response.
type Data interface {
	// Bind data into a structure
	Bind(interface{}) error
	// Raw data structure
	Raw() interface{}
}

// NewJSONData creates data object from a raw JSON message.
func NewJSONData(raw *json.RawMessage) Data {
	return &jsonData{raw: raw}
}

type jsonData struct {
	raw *json.RawMessage
}

func (r jsonData) Raw() (v interface{}) {
	//nolint:errcheck
	r.Bind(&v)
	return
}

func (r jsonData) Bind(v interface{}) error {
	if r.raw == nil {
		return nil
	}

	err := json.Unmarshal(*r.raw, v)
	if err != nil {
		return errors.NewInvalidParams(err.Error(), nil)
	}

	return nil
}

// NewStructData creates data object from a structure.
func NewStructData(raw interface{}) Data {
	return &structData{raw: raw}
}

type structData struct {
	raw interface{}
}

func (p *structData) Raw() interface{} {
	return p.raw
}

func (p *structData) Bind(v interface{}) error {
	if p.raw == nil {
		return nil
	}

	if reflect.TypeOf(v).Kind() != reflect.Ptr {
		return errors.NewInternal("Value used for payload should be a pointer", nil)
	}

	vp := reflect.ValueOf(p.raw)

	if reflect.TypeOf(p.raw).Kind() == reflect.Ptr {
		vp = vp.Elem()
	}

	vv := reflect.ValueOf(v).Elem()

	if vp.Type().ConvertibleTo(vv.Type()) {
		vp = vp.Convert(vv.Type())
	}

	if !vv.Type().AssignableTo(vp.Type()) {
		return errors.NewInvalidParams(
			fmt.Sprintf("Value of type %v is not assignable to value of type %v", vp.Type(), vv.Type()),
			nil,
		)
	}

	vv.Set(vp)

	return nil
}
