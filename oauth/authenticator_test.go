package oauth

import (
	"context"
	"reflect"
	"testing"
)

func TestAuthenticator_AuthenticateHTTP(t *testing.T) {
	srv := newOAuthServer(t)
	defer srv.Close()

	auth := NewAuthenticator(srv.URL)

	ctx, err := auth.AuthenticateHTTP(context.Background(), "Bearer", SomeToken)
	if err != nil {
		t.Error("Authentication suppose to be successful, got error instead:", err)
	}

	token, ok := TokenFromContext(ctx)
	if !ok {
		t.Error("Token is not present in context")
	}

	want := []string{"foo", "bar", "baz"}
	got := token.Scopes

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Scopes in the token do not match: want \"%v\", got \"%v\"", want, got)
	}
}