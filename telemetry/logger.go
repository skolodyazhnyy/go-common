package telemetry

type logger interface {
	Warning(msg string, args map[string]interface{})
}
