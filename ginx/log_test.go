package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type sliceLogger []string

func (l *sliceLogger) Error(msg string, data map[string]interface{}) {
	*l = append(*l, fmt.Sprintf("[ERROR] %v", msg))
}

func (l *sliceLogger) Debug(msg string, data map[string]interface{}) {
	*l = append(*l, fmt.Sprintf("[DEBUG] %v", msg))
}

func (l *sliceLogger) Flush() (all []string) {
	all = []string(*l)
	*l = nil
	return
}

func TestLog(t *testing.T) {
	log := &sliceLogger{}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(Log(log))
	router.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
	})
	router.GET("/error", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusInternalServerError)
	})

	srv := httptest.NewServer(router)
	defer srv.Close()

	t.Run("non 5xx are logged as debug", func(t *testing.T) {
		_, err := http.Get(srv.URL)
		if err != nil {
			t.Fatal("HTTP request to test server has failed:", err)
		}

		got := log.Flush()
		want := []string{`[DEBUG] GET request to / with status code 200`}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("Logged messages are incorrect, got: %v, want: %v", got, want)
		}
	})

	t.Run("5xx are logged as error", func(t *testing.T) {
		_, err := http.Get(srv.URL + "/error")
		if err != nil {
			t.Fatal("HTTP request to test server has failed:", err)
		}

		got := log.Flush()
		want := []string{`[ERROR] GET request to /error failed with status code 500`}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("Logged messages are incorrect, got: %v, want: %v", got, want)
		}
	})
}
