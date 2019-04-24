package telemetry

import "context"

type contextKey int

const contextTags contextKey = iota

// AppendContext with telemetry tags
func AppendContext(ctx context.Context, tags []string) context.Context {
	tags = append(TagsFromContext(ctx), tags...)
	return context.WithValue(ctx, contextTags, tags)
}

// TagsFromContext get tags from context
func TagsFromContext(ctx context.Context) []string {
	if ctx == nil {
		return nil
	}

	tags, _ := ctx.Value(contextTags).([]string)
	return tags
}
