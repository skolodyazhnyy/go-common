package httpx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type authenticatorFunc func(ctx context.Context, username, password string) (context.Context, bool)

func (f authenticatorFunc) AuthenticateHTTP(ctx context.Context, username, password string) (context.Context, bool) {
	return f(ctx, username, password)
}

func TestAuthenticate(t *testing.T) {
	auth := authenticatorFunc(func(ctx context.Context, kind, cred string) (context.Context, bool) {
		if kind == "Basic" && cred == "Z3Vlc3Q6Z3Vlc3Q=" {
			return context.WithValue(ctx, "token", "foobar"), true
		}

		return ctx, false
	})

	srv := httptest.NewServer(Authenticate(auth)(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if token, _ := req.Context().Value("token").(string); token != "foobar" {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}

		rw.WriteHeader(http.StatusOK)
	})))
	defer srv.Close()

	t.Run("valid credentials", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		if err != nil {
			t.Fatal("Request creation failed:", err)
		}

		req.SetBasicAuth("guest", "guest")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("Request to test server has failed:", err)
		}

		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode == 401 {
			t.Error("Test server replied 401, which probably means credentials are not properly parsed")
		}

		if resp.StatusCode != 200 {
			t.Errorf("Test server replied non-200, got %v instead", resp.StatusCode)
		}
	})

	t.Run("no credentials", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		if err != nil {
			t.Fatal("Request creation failed:", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("Request to test server has failed:", err)
		}

		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != 401 {
			t.Errorf("Test server replied non-401, got %v instead", resp.StatusCode)
		}
	})

	t.Run("invalid credentials", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, srv.URL, nil)
		if err != nil {
			t.Fatal("Request creation failed:", err)
		}

		req.SetBasicAuth("foo", "bar")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal("Request to test server has failed:", err)
		}

		defer resp.Body.Close() //nolint:errcheck

		if resp.StatusCode != 401 {
			t.Errorf("Test server replied non-401, got %v instead", resp.StatusCode)
		}
	})
}
