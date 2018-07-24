package httpx

import "net/http"

// Handlerx is an extended version of http.Handler, which adds a method to wrap handler into middleware
type Handlerx struct {
	http.Handler
}

// Handler creates extended version of http.Handler
func Handler(h http.Handler) *Handlerx {
	return &Handlerx{Handler: h}
}

// With adds middleware to handler
func (h *Handlerx) With(m ...func(http.Handler) http.Handler) *Handlerx {
	for _, x := range m {
		h.Handler = x(h.Handler)
	}

	return h
}
