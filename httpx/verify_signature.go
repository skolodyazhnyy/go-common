package httpx

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

// VerifySignature middleware verifies request's signature
func VerifySignature(sign signer, log logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Error("An error occurred while reading request body to verify signature", map[string]interface{}{
					"error":   err.Error(),
					"context": req.Context(),
				})
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			req.Body = ioutil.NopCloser(bytes.NewReader(body))

			header, err := sign.Sign(req.Header, body)
			if err != nil {
				log.Error("An error occurred while calculating request's signature", map[string]interface{}{
					"error":   err.Error(),
					"context": req.Context(),
				})
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			if req.Header.Get("X-Signature") != header {
				rw.WriteHeader(http.StatusForbidden)
				rw.Write([]byte("signature is invalid or missing")) //nolint:errcheck
				return
			}

			h.ServeHTTP(rw, req)
		})
	}
}
