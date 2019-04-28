package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"runtime"
)

type recovery interface {
	Error(msg string, data map[string]interface{})
}

// Recover middleware allows to gracefully handle panics raised during request processing
func Recover(log recovery) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error(fmt.Sprint(err), map[string]interface{}{"panic": location(5)})
				c.Writer.WriteHeader(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}

func location(skip int) map[string]interface{} {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(skip, fpcs)
	if n == 0 {
		return nil
	}

	fun := runtime.FuncForPC(fpcs[0] - 1)
	if fun == nil {
		return nil
	}

	file, line := fun.FileLine(fpcs[0] - 1)

	return map[string]interface{}{
		"function": path.Base(fun.Name()),
		"file":     path.Base(file),
		"line":     line,
	}
}
