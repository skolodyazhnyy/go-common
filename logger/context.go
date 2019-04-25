package logger

import "context"

type contextKey int

const recordContext contextKey = iota

// AppendContext with logging data
func AppendContext(ctx context.Context, r map[string]interface{}) context.Context {
	return context.WithValue(ctx, recordContext, merge(RecordFromContext(ctx), r))
}

// RecordFromContext fetches logging data from context
func RecordFromContext(ctx context.Context) map[string]interface{} {
	data, _ := ctx.Value(recordContext).(map[string]interface{})
	return data
}
