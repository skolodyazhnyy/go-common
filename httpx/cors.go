package httpx

import "net/http"

type cors interface {
	ValidateOrigin(origin string) bool
	PreflightHeaders() map[string][]string
	NormalHeaders() map[string][]string
}

func CORS(c cors) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if origin := req.Header.Get("Origin"); len(origin) != 0 {
				// is valid
				if !c.ValidateOrigin(origin) {
					rw.WriteHeader(http.StatusForbidden)
					return
				}

				// is pre-flight check
				if req.Method == http.MethodOptions {
					header := rw.Header()
					for k, v := range c.PreflightHeaders() {
						header[k] = v
					}

					rw.WriteHeader(http.StatusOK)
					return
				}

				// is normal request
				header := rw.Header()
				for k, v := range c.NormalHeaders() {
					header[k] = v
				}
			}

			h.ServeHTTP(rw, req)
		})
	}
}
