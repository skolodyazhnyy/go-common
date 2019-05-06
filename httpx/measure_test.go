package httpx_test

import (
	"context"
	"fmt"
	"github.com/magento-mcom/go-common/httpx"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type sliceMeter []string

func (m *sliceMeter) Timing(ctx context.Context, metric string, duration time.Duration, tags []string) {
	*m = append(*m, fmt.Sprintf("[Timing] %v %#v", metric, tags))
}

func (m *sliceMeter) Flush() (all []string) {
	all = []string(*m)
	*m = nil
	return
}

func TestMeasure(t *testing.T) {
	meter := &sliceMeter{}

	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(httpx.Measure(meter)(handler))
	defer srv.Close()

	_, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal("HTTP request to test server has failed:", err)
	}

	got := meter.Flush()
	want := []string{`[Timing] http_request []string{"method:GET", "status:200", "status_class:2xx"}`}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Reported metrics are incorrect, got: %v, want: %v", got, want)
	}
}
