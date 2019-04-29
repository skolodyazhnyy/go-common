package ginx

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type signer interface {
	Sign(map[string][]string, []byte) (string, error)
}

// VerifySignature middleware verifies request's signature
func VerifySignature(sign signer, log logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Error("An error occurred while reading request body to verify signature", map[string]interface{}{
				"error":   err.Error(),
				"context": c.Request.Context(),
			})
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))

		header, err := sign.Sign(c.Request.Header, body)
		if err != nil {
			log.Error("An error occurred while calculating request's signature", map[string]interface{}{
				"error":   err.Error(),
				"context": c.Request.Context(),
			})
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		if c.Request.Header.Get("X-Signature") != header {
			c.Writer.WriteHeader(http.StatusForbidden)
			c.Writer.Write([]byte("signature is invalid or missing")) //nolint:errcheck
			return
		}

		c.Next()
	}
}
