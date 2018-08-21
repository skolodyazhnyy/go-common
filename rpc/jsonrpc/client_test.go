package jsonrpc

import (
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"

	"net/http"

	"github.com/magento-mcom/go-common/rpc"
	"github.com/magento-mcom/go-common/rpc/errors"
)

type TestParams struct {
	A int
	B int
}

// Create JSONRPC client for endpoint "https://magento.com" using default HTTP client as transport.
func ExampleNewClient() {
	// create JSONRPC client
	cli := NewClient("https://magento.com", http.DefaultClient)

	// call "calc.add" with parameters
	res, err := cli.Call(rpc.NewRequest(nil, "calc.add", TestParams{A: 100, B: 20}))

	// report server error
	if rerr, ok := err.(errors.ErrorResponse); ok {
		fmt.Printf("Server replied with error: %v\n", rerr)
		return
	}

	// report rpc call failure
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}

	// read result value
	val := 0

	if err := res.Bind(&val); err != nil {
		fmt.Printf("Result is not an integer: %v\n", err)
		return
	}

	fmt.Printf("Result: %v\n", val)
}

// Create JSONRPC server using standard http server.
func ExampleNewServer() {
	// create JSONRPC server with some rpc.Handler
	srv := NewServer(rpc.HandlerFunc(func(req rpc.Request) (interface{}, error) {
		params := TestParams{}
		if err := req.Bind(&params); err != nil {
			return nil, errors.NewInvalidParams(err.Error(), nil)
		}

		return params.A + params.B, nil
	}))

	// start http server using JSONRPC Server as handler
	panic(http.ListenAndServe(":8080", srv))
}

func TestClient(t *testing.T) {
	srv := httptest.NewServer(NewServer(rpc.HandlerFunc(func(req rpc.Request) (interface{}, error) {
		switch req.Method {
		case "add":
			v := TestParams{}
			if err := req.Bind(&v); err != nil {
				return nil, err
			}

			return v.A + v.B, nil
		case "generic-error":
			return nil, fmt.Errorf("Generic error: %v", req.Method)
		default:
			return nil, errors.NewMethodNotFound(fmt.Sprintf("Method \"%v\" is not defined", req.Method), nil)
		}
	})))
	defer srv.Close()

	cli := NewClient(srv.URL, nil)

	t.Run("happy path", func(t *testing.T) {
		res, err := cli.Call(rpc.NewRequest(nil, "add", TestParams{A: 10, B: 81}))
		if err != nil {
			t.Fatal(err)
		}

		val := 0

		if err := res.Bind(&val); err != nil {
			t.Fatal(err)
		}

		if val != 91 {
			t.Errorf("Result is %v, but 91 is expected", val)
		}
	})

	t.Run("method not found", func(t *testing.T) {
		_, err := cli.Call(rpc.NewRequest(nil, "div", TestParams{}))
		if err == nil {
			t.Fatal("RPC request should fail, but it didn't")
		}

		rpcerr, ok := err.(errors.ErrorResponse)
		if !ok {
			t.Fatalf("Error type should be errors.ErrorResponse, but %v is received instead", reflect.TypeOf(err))
		}

		if rpcerr.Code != errors.CodeMethodNotFound {
			t.Fatalf("Error response code should be %v, but %v is received instead", errors.CodeMethodNotFound, rpcerr.Code)
		}
	})

	t.Run("generic error", func(t *testing.T) {
		_, err := cli.Call(rpc.NewRequest(nil, "generic-error", TestParams{}))
		if err == nil {
			t.Fatal("RPC request should fail, but it didn't")
		}

		rpcerr, ok := err.(errors.ErrorResponse)
		if !ok {
			t.Fatalf("Error type should be errors.ErrorResponse, but %v is received instead", reflect.TypeOf(err))
		}

		if rpcerr.Code != errors.CodeGeneric {
			t.Fatalf("Error response code should be %v, but %v is received instead", errors.CodeGeneric, rpcerr.Code)
		}
	})
}
