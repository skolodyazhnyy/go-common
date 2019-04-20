package httpx

import "net/http"

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rsk *responseWriter) WriteHeader(s int) {
	rsk.status = s
	rsk.ResponseWriter.WriteHeader(s)
}
