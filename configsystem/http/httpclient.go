package http

import "net/http"

type HttpClientInterface interface {
	Get(url string) (resp *http.Response, err error)
	Head(url string) (resp *http.Response, err error)
}
