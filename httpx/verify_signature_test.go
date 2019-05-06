package httpx_test

import (
	"bytes"
	"github.com/magento-mcom/go-common/httpx"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVerifySignature(t *testing.T) {
	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(httpx.VerifySignature(testSigner{}, &sliceLogger{})(handler))
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

		if resp.StatusCode == http.StatusForbidden {
			t.Error("Test server replied 401, which probably means signature is not well verified")
		}

		if resp.StatusCode != http.StatusOK {
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

		if resp.StatusCode == http.StatusOK {
			t.Error("Test server replied 200, which probably means signature is not well verified")
		}

		if resp.StatusCode != http.StatusForbidden {
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

		if resp.StatusCode == http.StatusOK {
			t.Error("Test server replied 200, which probably means signature is not well verified")
		}

		if resp.StatusCode != http.StatusForbidden {
			t.Errorf("Test server replied non-403, got %v instead", resp.StatusCode)
		}
	})
}
