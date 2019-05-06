# HTTPx

HTTPx package provides extensions for standard `net/http` package. 

## Usage

Package provides middleware (decorators) for HTTP client and HTTP server.

### Client

Client middleware can be incorporated using `httpx.NewClient` constructor. It takes http.Client as first argument and list of middleware to apply. 

```go
client := httpx.NewClient(
	&http.Client{Timeout: 30 * time.Second},
	httpx.WithHeaders(...),
	httpx.WithCredentials(...),	
)
```

#### WithHeaders

WithHeaders middleware allows to set additional HTTP headers on every request.

```go
client := httpx.NewClient(http.DefaultClient, httpx.WithHeaders(map[string][]string{
	"User-Agent": {"Go HTTPx Client"},
	"X-Powered-By": {"Go"},
}))
```

#### WithCredentials

WithCredentials middleware allows to set authorization header. It's enabled only if server replies 401 to a request. 

When middleware detects 401 response it uses credentials provider to generate new token and repeats HTTP request again with this token. In case server still replies 401, it returns response as it is. 

Once token is generated it is used for every subsequent request. 

```go
client := httpx.NewClient(http.DefaultClient, httpx.WithCredentials(
	oauth.NewClientCredentials(...),
))
```

#### WithSignature

WithSignature middleware allows to add signature header to every request. It calculates signature and adds `X-Signature` header to requests.

```go
client := httpx.NewClient(http.DefaultClient, httpx.WithSignature(
	sha1signer.New(),
))
```

### Server

Server middleware can be incorporated using [gorilla/mux](https://github.com/gorilla/mux) package, using `Use` method of the router.

```go
router := mux.New()
router.Use(
	httpx.Authenticate(...)
	httpx.Measure(...)
	httpx.Recover(...)
	httpx.Log(...)
)

```

#### Authenticate

Authenticate middleware abstracts authentication for every request. It takes authenticator instance which adds authentication information into request's context. 

```go
router.Use(httpx.Authenticate(oauth.NewAuthenticator(...), log.Channel("http")))
```

Later request can be authorized using authentication information from request's context. See `oauth` package documentation for authorization example.

#### Log

Log middleware allows to log every request received by HTTP server.

```go
router.Use(httpx.Log(logger))
```

#### Recover

Recover middleware allows to recover in case HTTP handler rises panic. It takes logger to log error in case of panic.

```go
router.Use(httpx.Recover(logger))
```

#### Measure

Measure middleware allows to add metrics for HTTP server. It reports number of requests processed by server and their latency.

```go
router.Use(httpx.Measure(telemetry))
```

#### VerifySignature

VerifySignature middleware allows to check request's signature in X-Signature header and reject request if it does not match expected value.

```go
router.Use(httpx.VerifySignature(signer, log))
```
