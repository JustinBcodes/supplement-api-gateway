package userctx

import "context"

type key int

const (
	keySub key = iota
)

func WithSub(ctx context.Context, sub string) context.Context {
	return context.WithValue(ctx, keySub, sub)
}

func Sub(ctx context.Context) string {
	if v, ok := ctx.Value(keySub).(string); ok {
		return v
	}
	return ""
}


