package httpx_test

import (
	"errors"
	"github.com/magento-mcom/go-common/httpx"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SignerFunc func(map[string][]string, []byte) (string, error)

func (f SignerFunc) Sign(headers map[string][]string, body []byte) (string, error) {
	return f(headers, body)
}

func TestWithSignature(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-Signature") != "FooSignature" {
			t.Error("Request is not signed")
		}

		rw.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	t.Run("signature-ok", func(t *testing.T) {
		cli := httpx.NewClient(
			&http.Client{Timeout: time.Second},
			httpx.WithSignature(SignerFunc(func(headers map[string][]string, body []byte) (string, error) {
				return "FooSignature", nil
			})),
		)

		req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		if err != nil {
			t.Fatal("Request can not be created:", err)
		}

		resp, err := cli.Do(req)
		if err != nil {
			t.Fatal("Request to test server has failed:", err)
		}

		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != 200 {
			t.Errorf("Response status code is not 200, got %v instead", resp.StatusCode)
		}
	})

	t.Run("signature-failed", func(t *testing.T) {
		want := errors.New("signer failed")

		cli := httpx.NewClient(
			&http.Client{Timeout: time.Second},
			httpx.WithSignature(SignerFunc(func(headers map[string][]string, body []byte) (string, error) {
				return "", want
			})),
		)

		req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		if err != nil {
			t.Fatal("Request can not be created:", err)
		}

		_, err = cli.Do(req)
		if err == nil {
			t.Fatal("Request suppose to fail but it didn't")
		}

		if err != want {
			t.Errorf("Error does not match expected value: want %v, got %v", err, want)
		}
	})
}
