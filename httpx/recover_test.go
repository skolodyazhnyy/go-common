package httpx_test

import (
	"errors"
	"github.com/magento-mcom/go-common/httpx"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRecovery(t *testing.T) {
	log := &sliceLogger{}

	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		panic(errors.New("oopsie daisy"))
	})

	srv := httptest.NewServer(httpx.Recover(log)(handler))
	defer srv.Close()

	t.Run("panics are logged as errors", func(t *testing.T) {
		_, err := http.Get(srv.URL)
		if err != nil {
			t.Fatal("HTTP request to test server has failed:", err)
		}

		got := log.Flush()
		want := []string{`[ERROR] oopsie daisy`}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("Logged messages are incorrect, got: %v, want: %v", got, want)
		}
	})
}
