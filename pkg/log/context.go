package log

import "context"

type key int

const logContextKey key = iota

// FromContext 返回 ctx 上 log 键的值。
func FromContext(ctx context.Context) Logger {
	if ctx != nil {
		logger := ctx.Value(logContextKey)
		if logger != nil {
			return logger.(Logger)
		}
	}
	return WithName("Unknown-Context")
}
