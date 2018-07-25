package log

import "context"

type contextKey int

const recordContext contextKey = iota

// AppendContext with logging data
func AppendContext(ctx context.Context, r map[string]interface{}) context.Context {
	return context.WithValue(ctx, recordContext, RecordFromContext(ctx).With(r))
}

// RecordFromContext fetches logging data from context
func RecordFromContext(ctx context.Context) R {
	data, _ := ctx.Value(recordContext).(R)
	return data
}
