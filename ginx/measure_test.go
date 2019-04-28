package ginx

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
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

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(Measure(meter))
	router.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(router)
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
