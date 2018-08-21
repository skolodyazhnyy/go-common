package rpc

import "net/http"

// Client provides interface to make calls to remote server.
// It returns a data container with result of RPC call or an error in case
// RPC request failed, or server returned an error.
type Client interface {
	Call(Request) (Data, error)
}

// Transport provides generic interface for HTTP based RPC clients.
type Transport interface {
	Do(*http.Request) (*http.Response, error)
}

// TransportMiddleware extends transport with additional capabilities.
// With transport middleware you can simply set an additional header, or perform
// more complex tasks like fetch authorization token, or adding caching layer.
type TransportMiddleware func(Transport) Transport
