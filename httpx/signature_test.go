package httpx_test

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/magento-mcom/go-common/httpx"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testSigner struct {
	failure error
}

func (s testSigner) Sign(headers map[string][]string, body []byte) (string, error) {
	h := md5.New()
	h.Write(body) //nolint:errcheck
	return fmt.Sprintf("%x", h.Sum(nil)), s.failure
}

func TestWithSignature(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if sign := req.Header.Get("X-Signature"); sign != "d41d8cd98f00b204e9800998ecf8427e" {
			t.Error("Request is not signed or signature is invalid, got:", sign)
		}

		rw.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	t.Run("signature-ok", func(t *testing.T) {
		cli := httpx.NewClient(
			&http.Client{Timeout: time.Second},
			httpx.WithSignature(testSigner{}),
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

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Response status code is not 200, got %v instead", resp.StatusCode)
		}
	})

	t.Run("signature-failed", func(t *testing.T) {
		want := errors.New("signer failed")

		cli := httpx.NewClient(
			&http.Client{Timeout: time.Second},
			httpx.WithSignature(testSigner{failure: want}),
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
