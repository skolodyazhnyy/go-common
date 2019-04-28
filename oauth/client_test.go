package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const (
	SomeClientID     = "red"
	SomeClientSecret = "green"
	SomeToken        = "yellow"
	InvalidToken     = "foo"
)

func newOAuthServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/token/" + SomeToken + "/scopes":
			if req.Method != http.MethodGet {
				rw.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			rw.WriteHeader(http.StatusOK)
			//nolint:errcheck
			rw.Write([]byte(`{"access_token": "` + SomeToken + `", "scopes": ["foo", "bar", "baz"]}`))
		case "/token/" + InvalidToken + "/scopes":
			if req.Method != http.MethodGet {
				rw.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			rw.WriteHeader(http.StatusUnauthorized)
		case "/oauth/token":
			if req.Method != http.MethodPost {
				rw.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			if err := req.ParseForm(); err != nil {
				t.Error("OAuth server can not parse request:", err)
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			if form := req.PostForm; form.Get("grant_type") != "client_credentials" || form.Get("client_id") != SomeClientID || form.Get("client_secret") != SomeClientSecret {
				rw.WriteHeader(http.StatusUnauthorized)
				return
			}

			rw.WriteHeader(http.StatusOK)
			//nolint:errcheck
			rw.Write([]byte(`{"access_token": "` + SomeToken + `"}`))
		default:
			t.Errorf("Invalid OAuth request path: %v", req.URL.Path)
		}
	}))
}

func TestClient_ClientCredentials(t *testing.T) {
	srv := newOAuthServer(t)
	defer srv.Close()

	cli := NewClient(srv.URL)

	t.Run("valid-credentials", func(t *testing.T) {
		token, err := cli.ClientCredentials(context.Background(), SomeClientID, SomeClientSecret)
		if err != nil {
			t.Fatal("Client returned an error:", err)
		}

		if token != SomeToken {
			t.Fatalf("Token does not match: want %v, got %v", SomeToken, token)
		}
	})

	t.Run("invalid-credentials", func(t *testing.T) {
		_, err := cli.ClientCredentials(context.Background(), "foo", "bar")
		if err == nil {
			t.Fatal("Client should have returned an error")
		}

		if err != ErrInvalidCredentials {
			t.Fatalf("Client should have returned ErrInvalidCredentials, but it returned: %v instead", err)
		}
	})
}

func TestClient_Scopes(t *testing.T) {
	srv := newOAuthServer(t)
	defer srv.Close()

	cli := NewClient(srv.URL)

	t.Run("valid-token", func(t *testing.T) {
		got, err := cli.Scopes(context.Background(), SomeToken)
		if err != nil {
			t.Fatal("Client returned an error:", err)
		}

		if want := []string{"foo", "bar", "baz"}; !reflect.DeepEqual(got, want) {
			t.Fatalf("Scopes do not match: want %v, got %v", want, got)
		}
	})

	t.Run("invalid-token", func(t *testing.T) {
		_, err := cli.Scopes(context.Background(), InvalidToken)
		if err == nil {
			t.Fatal("Client should have returned an error")
		}

		if err != ErrInvalidToken {
			t.Fatalf("Client should have returned ErrTokenInvalid, but it returned: %v instead", err)
		}
	})
}
