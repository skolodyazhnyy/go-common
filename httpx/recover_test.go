package httpx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRecovery(t *testing.T) {
	log := &sliceLogger{}
	
	srv := httptest.NewServer(Recover(log)(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		panic(errors.New("oopsie daisy"))
	})))
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
