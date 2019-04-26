package oauth

// Token represents OAuth token
type Token struct {
	Scopes []string
}

// Has required scope
func (t Token) Has(want string) bool {
	for _, scope := range t.Scopes {
		if scope == want {
			return true
		}
	}

	return false
}
