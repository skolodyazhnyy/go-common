package httpx

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

type signer interface {
	Sign(map[string][]string, []byte) (string, error)
}

// WithSignature signs outgoing HTTP call
// This middleware adds X-Signature header with request hash-sum
func WithSignature(sign signer) func(Client) Client {
	return func(c Client) Client {
		return ClientFunc(func(req *http.Request) (resp *http.Response, err error) {
			var body []byte

			if req.Body != nil {
				if body, err = ioutil.ReadAll(req.Body); err != nil {
					return nil, err
				}

				req.Body = ioutil.NopCloser(bytes.NewReader(body))
			}

			header, err := sign.Sign(req.Header, body)
			if err != nil {
				return nil, err
			}

			req.Header.Set("X-Signature", header)

			return c.Do(req)
		})
	}
}
