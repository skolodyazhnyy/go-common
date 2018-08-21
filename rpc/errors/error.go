package errors

// Error codes are taken from JSONRPC/XMLRPC specifications.
//
// Error code prefix types:
// 	-310: Service Bus error (i.e. 4xx http error codes)
// 	-311: Proxy error (error in service bus while proxying a request)
// 	-320: Server-level error
// 	-325: Application-level error
// 	-326: Given request is invalid
// 	-327: Error during request parsing
const (
	CodeProbeFailed      int = -31001 // Service bus: application probing failed
	CodeGeneric          int = -32000 // Generic error code
	CodeEndpointDisabled int = -32501 // Endpoint is disabled
	CodeParse            int = -32700 // Error during request parsing
	CodeInvalidRequest   int = -32600 // Generic invalid request error
	CodeMethodNotFound   int = -32601 // Invalid request: method not found
	CodeInvalidParams    int = -32602 // Invalid request: method parameters are invalid
	CodeInternal         int = -32603 // Invalid request: internal error
	CodeAccessDenied     int = -32604 // Invalid request: access denied
)

// ShouldRetry tells if given error code should be retried
func ShouldRetry(err error) bool {
	e, ok := err.(ErrorResponse)
	if !ok { // error is not bus.Error, probably some unexpected error, lets retry
		return true
	}

	switch e.Code {
	case 0, // zero normally would mean there was no error code specified, so we retry it just in case
		CodeEndpointDisabled,      // endpoint is currently disabled but we will retry later if delegated call
		CodeInternal, CodeGeneric: // internal and generic codes are for unhandled exceptions on the server
		return true
	default:
		return false
	}
}

// ErrorResponse represents a confirmed response from a remote endpoint
type ErrorResponse struct {
	msg  string
	Code int
	Data interface{}
}

func (be ErrorResponse) Error() string {
	return be.msg
}

// DecodingError represents a response which could not be properly decoded and thus it's considered a communication failure
type DecodingError struct {
	msg         string
	RawResponse []byte
}

func (bd DecodingError) Error() string {
	return bd.msg
}

// HTTPResponseError is error when HTTP response status code is not 2xx
type HTTPResponseError struct {
	msg    string
	Status int
}

func (e HTTPResponseError) Error() string {
	return e.msg
}

// ServiceNotFoundError is used to signal non existing services and its used to return a 404 response back
type ServiceNotFoundError struct {
	msg string
}

func (se ServiceNotFoundError) Error() string {
	return se.msg
}

// New constructs a new error with custom error code and message
func New(code int, msg string, data interface{}) ErrorResponse {
	return ErrorResponse{msg: msg, Code: code, Data: data}
}

// NewParse constructs a new request parsing error
func NewParse(msg string, data interface{}) ErrorResponse {
	return New(CodeParse, msg, data)
}

// NewInvalidRequest constructs a new invalid request error
func NewInvalidRequest(msg string, data interface{}) ErrorResponse {
	return New(CodeInvalidRequest, msg, data)
}

// NewMethodNotFound constructs a new method not found error
func NewMethodNotFound(msg string, data interface{}) ErrorResponse {
	return New(CodeMethodNotFound, msg, data)
}

// NewServiceNotFound constructs a service not found error
func NewServiceNotFound(msg string) error {
	return ServiceNotFoundError{msg: msg}
}

// NewEndpointDisabled constructs a new endpoit is disabled error
func NewEndpointDisabled(msg string, data interface{}) ErrorResponse {
	return New(CodeEndpointDisabled, msg, data)
}

// NewInvalidParams constructs a new invalid parameters error
func NewInvalidParams(msg string, data interface{}) ErrorResponse {
	return New(CodeInvalidParams, msg, data)
}

// NewInternal constructs a new invalid request internal error
func NewInternal(msg string, data interface{}) ErrorResponse {
	return New(CodeInternal, msg, data)
}

// NewDecodingResponseError constructs a new decoding error
func NewDecodingResponseError(msg string, rawResponse []byte) DecodingError {
	return DecodingError{
		msg:         msg,
		RawResponse: rawResponse,
	}
}

// NewHTTPResponseError constructs a new HTTP response error
func NewHTTPResponseError(msg string, status int) HTTPResponseError {
	return HTTPResponseError{
		msg:    msg,
		Status: status,
	}
}
