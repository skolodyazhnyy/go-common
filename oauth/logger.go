package oauth

type logger interface {
	Error(string, map[string]interface{})
}

type nopLogger struct {
}

func (nopLogger) Error(string, map[string]interface{}) {
}
