# OAuth

OAuth package provides OAuth client (compatible with MCOM OAuth server) and authentication services for HTTP middleware.

## Usage

Package provides OAuth client and services for HTTP server and client middleware.

### Client

OAuth client which can be used on its own to create and inspect tokens.

```go
client := oauth.NewClient("https://auth.mcom.magento.com")
```

Create token using client credentials

```go
token, err := client.ClientCredentials(ctx, "client-id", "client-secret")
if err == oauth.ErrInvalidCredentials { // if needed invalid credentials error can be handled differently
	panic(err)
}

if err != nil { // generic error
	panic(err)
}

fmt.Printf("Token: %#v\n", token)
```

Client allows to fetch token's scopes

```go
scopes, err := client.Scopes(ctx, token)
if err != oauth.ErrInvalidToken { // if needed invalid token error can be handled differently
	panic(err)
}

if err != nil { // generic error
	panic(err)
}

fmt.Printf("Scopes: %#v\n", scopes)
```

### HTTP server

Package implements Authenticator, compatible with httpx server authentication middleware:

```go
router := mux.New()
router.Use(httpx.Authenticate(oauth.NewAuthenticator("https://auth.mcom.magento.com")))
``` 

This middleware would add `oauth.Token` into request's context, which can be later used in `http.Handler` or `rpc.Handler` to perform authorization

```go
token, ok := oauth.TokenFromContext(req.Context())
if !ok { // check if token is set
	rw.WriteHeader(http.StatusUnauthenticated)
	return
}

if !token.Has("user_admin") { // check if token has scope
	rw.WriteHeader(http.StatusForbidden)
	return	
}
```  

### HTTP client

Package implements credentials provider, compatible with httpx credentials middleware:

```go
client := httpx.NewClient(
    &http.Client{Timeout: 30 * time.Second},
    httpx.WithCredentials(oauth.NewClientCredentials("https://oauth.magento.com", "mcom", "mc0m-s3cr3t")),
)

client.Do(req)
```

Middleware will automatically add token to HTTP requests made with httpx client and refresh it in case server responds with 401.