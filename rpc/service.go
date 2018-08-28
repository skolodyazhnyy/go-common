package rpc

import (
	"fmt"

	"github.com/magento-mcom/go-common/rpc/errors"
)

// ServiceDesc describes a service
type ServiceDesc struct {
	ServiceName string
	Methods     []MethodDesc
}

// MethodDesc describes a method of the service
type MethodDesc struct {
	MethodName string
	Handler    HandlerFunc
}

// NewService constructs rpc.Handler which handles request using service descriptor.
// Service descriptor provide mapping between method and method handler. It helps to
// describe API and does RPC request routing based on the method name.
func NewService(descriptors ...ServiceDesc) Handler {
	handlers := map[string]HandlerFunc{}

	for _, d := range descriptors {
		for _, m := range d.Methods {
			handlers[d.ServiceName+"."+m.MethodName] = m.Handler
		}
	}

	return HandlerFunc(func(req Request) (interface{}, error) {
		h, ok := handlers[req.Method]
		if !ok {
			return nil, errors.NewMethodNotFound(fmt.Sprintf("Method \"%v\" is not defined", req.Method), nil)
		}

		return h(req)
	})
}
