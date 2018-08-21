package rpc

import (
	"context"
	"encoding/json"
)

// NewRequest constructs an RPC request object.
//
// It accepts any data structure as container for parameters. Although,
// if params implements Data interface, it will be used to marshal or bind
// parameters.
//
// If context is not provided, background context will be used.
func NewRequest(ctx context.Context, method string, params interface{}) Request {
	if ctx == nil {
		ctx = context.Background()
	}

	p, ok := params.(Data)
	if !ok {
		p = NewStructData(params)
	}

	return Request{
		Method:  method,
		Context: ctx,
		Params:  p,
	}
}

// NewJSONRequest constructs an RPC request object for JSON encoded parameters.
// This constructor is useful when constructing Request from JSON.
func NewJSONRequest(ctx context.Context, method string, params *json.RawMessage) Request {
	return NewRequest(ctx, method, NewJSONData(params))
}

// Request represents RPC call.
type Request struct {
	Method  string
	Params  Data
	Context context.Context
}

// Raw returns underlying structure which holds parameters.
func (r *Request) Raw() interface{} {
	if r.Params == nil {
		return nil
	}

	return r.Params.Raw()
}

// Bind parameters to a structure.
func (r *Request) Bind(v interface{}) error {
	if r.Params == nil {
		return nil
	}

	return r.Params.Bind(v)
}

// Value fetches a value from request context.
func (r *Request) Value(key interface{}) interface{} {
	if r.Context == nil {
		return nil
	}

	return r.Context.Value(key)
}

// String makes string representation of the request.
func (r *Request) String() string {
	return r.Method
}
