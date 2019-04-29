package ginx_test

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/magento-mcom/go-common/ginx"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testSigner struct {
	failure error
}

func (s testSigner) Sign(headers map[string][]string, body []byte) (string, error) {
	h := md5.New()
	h.Write(body) //nolint:errcheck
	return fmt.Sprintf("%x", h.Sum(nil)), s.failure
}

func TestVerifySignature(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(ginx.VerifySignature(testSigner{}, &sliceLogger{}))
	router.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(router)
	defer srv.Close()

	t.Run("valid signature", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL, bytes.NewReader([]byte("hello world")))
		if err != nil {
			t.Fatal("Request creation failed:", err)
		}

		req.Header.Set("X-Signature", "5eb63bbbe01eeed093cb22bb8f5acdc3")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("Request to test server has failed:", err)
		}

		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode == 403 {
			t.Error("Test server replied 401, which probably means signature is not well verified")
		}

		if resp.StatusCode != 200 {
			t.Errorf("Test server replied non-200, got %v instead", resp.StatusCode)
		}
	})

	t.Run("invalid signature", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL, bytes.NewReader([]byte("hello world")))
		if err != nil {
			t.Fatal("Request creation failed:", err)
		}

		req.Header.Set("X-Signature", "hello world")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("Request to test server has failed:", err)
		}

		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode == 200 {
			t.Error("Test server replied 200, which probably means signature is not well verified")
		}

		if resp.StatusCode != 403 {
			t.Errorf("Test server replied non-403, got %v instead", resp.StatusCode)
		}
	})

	t.Run("no signature", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL, bytes.NewReader([]byte("hello world")))
		if err != nil {
			t.Fatal("Request creation failed:", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("Request to test server has failed:", err)
		}

		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode == 200 {
			t.Error("Test server replied 200, which probably means signature is not well verified")
		}

		if resp.StatusCode != 403 {
			t.Errorf("Test server replied non-403, got %v instead", resp.StatusCode)
		}
	})
}
