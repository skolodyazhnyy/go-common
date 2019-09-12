# Signature

Signature package provides services for signing HTTP requests. It meant to be used with [`httpx.WithSignature`](../httpx/README.md#withsignature), [`httpx.VerifySignature`](../httpx/README.md#verifysignature) and  [`ginx.VerifySignature`](../ginx/README.md#verifysignature) middleware.

## Usage

With HTTP server:

```golang
router.Use(httpx.VerifySignature(
    signature.SHA1WithSecret("secret"),
    log,
))
``` 

With Gin server:

```golang
router.Use(ginx.VerifySignature(
    signature.SHA1WithSecret("secret"),
    log,
))
``` 

With HTTP client:

```golang
client := httpx.NewClient(http.DefaultClient, httpx.WithSignature(
	signature.SHA1WithSecret("secret"),
))
``` 
