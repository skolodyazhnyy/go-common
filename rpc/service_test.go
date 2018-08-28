package rpc

import (
	"fmt"
	"testing"
)

type ExampleParams struct {
	A int
	B int
}

// Creating RPC handler using service descriptor. The service will provide two RPC
// methods "calc.add" and "calc.sub".
func ExampleNewService() {
	// define service
	handler := NewService(ServiceDesc{
		ServiceName: "calc",
		Methods: []MethodDesc{
			{
				MethodName: "add",
				Handler: func(req Request) (interface{}, error) {
					params := ExampleParams{}
					if err := req.Bind(&params); err != nil {
						return nil, err
					}

					return params.A + params.B, nil
				},
			},
			{
				MethodName: "sub",
				Handler: func(req Request) (interface{}, error) {
					params := ExampleParams{}
					if err := req.Bind(&params); err != nil {
						return nil, err
					}

					return params.A - params.B, nil
				},
			},
		},
	})

	r, _ := handler.Handle(NewRequest(nil, "calc.add", ExampleParams{A: 10, B: 2}))
	fmt.Printf("Result: %v\n", r)

	_, err := handler.Handle(NewRequest(nil, "calc.div", ExampleParams{A: 10, B: 2}))
	fmt.Printf("Error: %v\n", err)

	// Output:
	// Result: 12
	// Error: Method "calc.div" is not defined
}

func TestService_Handle(t *testing.T) {
	ctrl := NewService(ServiceDesc{
		ServiceName: "foo",
		Methods: []MethodDesc{
			{
				MethodName: "bar",
				Handler: func(req Request) (interface{}, error) {
					return req.Raw().(int) + 1, nil
				},
			},
		},
	})

	t.Run("Handle existent method", func(t *testing.T) {
		v, _ := ctrl.Handle(NewRequest(nil, "foo.bar", 100))
		if v == nil || v.(int) != 101 {
			t.Error("Handler returned wrong result")
		}
	})

	t.Run("Handle non-existent method", func(t *testing.T) {
		_, e := ctrl.Handle(NewRequest(nil, "bazz", nil))
		if e == nil {
			t.Error("Handler should return an error when handling invalid method")
		}
	})
}
