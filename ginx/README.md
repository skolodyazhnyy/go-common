# Ginx

Ginx package provides extensions for [gin-gonic](https://github.com/gin-gonic/gin) HTTP web framework. 

## Usage

Server middleware can be incorporated using `Use` method of the router.

```go
router := gin.New()
router.Use(
	ginx.Authenticate(...)
	ginx.Measure(...)
	ginx.Recover(...)
	ginx.Log(...)
)

```

### Authenticate

Authenticate middleware abstracts authentication for every request. It takes authenticator instance which adds authentication information into request's context. 

```go
router.Use(ginx.Authenticate(
	oauth.Authenticator(...),
))
```

Later request can be authorized using authentication information from request's context. See `oauth` package documentation for authorization example.

### Log

Log middleware allows to log every request received by HTTP server.

```go
router.Use(ginx.Log(logger))
```

### Recover

Recover middleware allows to recover in case HTTP handler rises panic. It takes logger to log error in case of panic.

```go
router.Use(ginx.Recover(logger))
```

### Measure

Measure middleware allows to add metrics for HTTP server. It reports number of requests processed by server and their latency.

```go
router.Use(ginx.Measure(telemetry))
```

### VerifySignature

VerifySignature middleware allows to check request's signature in X-Signature header and reject request if it does not match expected value.

```go
router.Use(ginx.VerifySignature(signer, log))
```
