/*
Copyright © 2025 lixw
*/
package safego

import (
	"context"
	"log/slog"
	"runtime/debug"
)

func Go(ctx context.Context, fn func(context.Context)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				slog.ErrorContext(ctx, "recovered from panic", "err", r, "stack", string(debug.Stack()))
			}
		}()

		fn(ctx)
	}()
}
