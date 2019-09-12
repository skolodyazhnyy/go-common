package signature

import "testing"

func TestSHA1WithSecret(t *testing.T) {
	tests := []struct {
		name      string
		secret    string
		body      []byte
		signature string
	}{
		{
			name:      "simple signature",
			secret:    "github-token",
			body:      []byte("magento"),
			signature: "sha1=931a628da747c612f212f0ba7cb0fc6746b723f3",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sgn := SHA1WithSecret(test.secret)

			got, err := sgn.Sign(nil, test.body)
			if err != nil {
				t.Fatal("Signer has returned an error:", err)
			}

			if want := test.signature; got != want {
				t.Errorf("Signature does not match, want: %#v, got %#v", want, got)
			}
		})
	}
}
