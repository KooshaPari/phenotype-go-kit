package logctx_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/KooshaPari/phenotype-go-kit/logctx"
)

func TestWithLoggerAndFrom(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	ctx := logctx.WithLogger(context.Background(), logger)
	got := logctx.From(ctx)
	if got != logger {
		t.Fatal("expected same logger back from context")
	}
}

func TestFromPanicsWithoutLogger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic when no logger in context")
		}
	}()
	logctx.From(context.Background())
}
