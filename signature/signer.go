package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
)

// Signer provides interface to sign a request
type Signer interface {
	// Sign method returns signature value for given headers and body
	Sign(map[string][]string, []byte) (string, error)
}

// SignerFunc implements Signer interface
type SignerFunc func(map[string][]string, []byte) (string, error)

// Sign method returns signature value for given headers and body
func (f SignerFunc) Sign(headers map[string][]string, body []byte) (string, error) {
	return f(headers, body)
}

// SHA1WithSecret signer
func SHA1WithSecret(secret string) Signer {
	return SignerFunc(func(headers map[string][]string, body []byte) (string, error) {
		mac := hmac.New(sha1.New, []byte(secret))

		_, _ = mac.Write(body)

		return "sha1=" + hex.EncodeToString(mac.Sum(nil)), nil
	})
}
