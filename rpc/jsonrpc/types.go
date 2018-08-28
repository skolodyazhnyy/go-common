package jsonrpc

import "encoding/json"

// ContextRequestID is request ID used in JSONRPC envelop.
const ContextRequestID contextKey = "request-id"

const version = "2.0"

type contextKey string

type requestDataEncode struct {
	Version string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type requestDataDecode struct {
	Version string           `json:"jsonrpc"`
	ID      interface{}      `json:"id"`
	Method  string           `json:"method"`
	Params  *json.RawMessage `json:"params"`
}

type responseData struct {
	Version string           `json:"jsonrpc"`
	ID      interface{}      `json:"id"`
	Result  *json.RawMessage `json:"result"`
	Error   *errorData       `json:"error"`
}

type resultResponseData struct {
	Version string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result"`
}

type errorResponseData struct {
	Version string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Error   *errorData  `json:"error"`
}

type errorData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
