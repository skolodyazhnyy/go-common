package jsonrpc

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/magento-mcom/go-common/rpc"
	"github.com/magento-mcom/go-common/rpc/errors"
	"gopkg.in/gin-gonic/gin.v1"
)

// GinServer constructs new JSONRPC HTTP handler for gin-tonic http package.
func GinServer(handler rpc.Handler, mw ...rpc.Middleware) func(ctxt *gin.Context) {
	h := NewServer(handler, mw...)

	return func(ctx *gin.Context) {
		h(ctx.Writer, ctx.Request)
	}
}

// NewServer constructs new JSONRPC HTTP handler for standard http package.
func NewServer(handler rpc.Handler, mw ...rpc.Middleware) http.HandlerFunc {
	handler = rpc.Decorate(handler, mw...)

	return func(rw http.ResponseWriter, r *http.Request) {
		req, resp, err := handle(r, handler)

		if _, ok := err.(errors.ServiceNotFoundError); ok {
			rw.WriteHeader(http.StatusNotFound)
			rw.Write([]byte("404 page not found"))
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)

		if err != nil {
			writeError(rw, req, err)
		} else {
			writeResponse(rw, req, resp)
		}
	}
}

func handle(r *http.Request, handler rpc.Handler) (rpc.Request, interface{}, error) {
	req, err := readRequest(r.Context(), r.Body)
	if err != nil {
		return req, nil, err
	}

	resp, err := handler.Handle(req)

	return req, resp, err
}

func readRequest(ctx context.Context, in io.Reader) (rpc.Request, error) {
	data := requestDataDecode{}

	b, err := ioutil.ReadAll(in)
	if err != nil {
		return rpc.Request{}, err
	}

	if err := json.Unmarshal(b, &data); err != nil {
		return rpc.Request{}, errors.NewParse("Parse error: "+err.Error(), string(b))
	}

	if data.Version != version {
		return rpc.Request{}, errors.NewInvalidRequest("JSONRPC version should be exactly 2.0", nil)
	}

	if data.Method == "" {
		return rpc.Request{}, errors.NewInvalidRequest("RPC method name can not be empty", nil)
	}

	ctx = context.WithValue(ctx, ContextRequestID, data.ID)

	return rpc.NewJSONRequest(ctx, data.Method, data.Params), nil
}

func writeError(w io.Writer, req rpc.Request, e error) error {
	errData := &errorData{
		Code:    errors.CodeGeneric,
		Message: e.Error(),
		Data:    nil,
	}

	if err, ok := e.(errors.ErrorResponse); ok {
		errData.Code = err.Code
		errData.Data = err.Data
	}

	if err, ok := e.(errors.DecodingError); ok {
		errData.Data = string(err.RawResponse)
	}

	data, err := json.Marshal(errorResponseData{
		Version: version,
		ID:      req.Value(ContextRequestID),
		Error:   errData,
	})

	if err != nil {
		return err
	}

	w.Write(data)

	return nil
}

func writeResponse(w io.Writer, req rpc.Request, resp interface{}) error {
	data, err := json.Marshal(resultResponseData{
		Version: version,
		ID:      req.Value(ContextRequestID),
		Result:  resp,
	})

	if err != nil {
		return err
	}

	w.Write(data)

	return nil
}
