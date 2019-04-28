package ginx

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestRecovery(t *testing.T) {
	log := &sliceLogger{}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(Recovery(log))
	router.GET("/", func(c *gin.Context) {
		panic(errors.New("oopsie daisy"))
	})

	srv := httptest.NewServer(router)
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
