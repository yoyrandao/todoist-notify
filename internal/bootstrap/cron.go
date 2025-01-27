package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func Schedule(ctx context.Context, fn func(context.Context) error, name string, interval time.Duration) {
	go func() {
		runInternal(ctx, fn, name, interval)

		for {
			select {
			case <-ctx.Done():
				return

			case <-time.After(interval):
				runInternal(ctx, fn, name, interval)
			}
		}
	}()
}

func runInternal(ctx context.Context, fn func(context.Context) error, name string, interval time.Duration) {
	if err := fn(ctx); err != nil {
		slog.Error("", "error", err.Error())
	}

	slog.Debug(fmt.Sprintf("job %s successed, next run in %f seconds", name, interval.Seconds()))
}
