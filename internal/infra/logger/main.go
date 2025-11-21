package logger

import (
	"context"
	"log/slog"

	sloghttp "github.com/samber/slog-http"
)

type Logger struct {
	base *slog.Logger
}

type LoggerInterface interface {
	Info(ctx context.Context, msg string, attrs ...slog.Attr)
	Error(ctx context.Context, msg string, attrs ...slog.Attr)
}

func NewLogger(base *slog.Logger) *Logger {
	return &Logger{base: base}
}

func grabRequestID(ctx context.Context, attrs []slog.Attr) []any {
	if reqID := sloghttp.GetRequestIDFromContext(ctx); reqID != "" {
		attrs = append(attrs, slog.String("id", reqID))
	}
	// convert []slog.Attr to []any
	args := make([]any, len(attrs))
	for i, a := range attrs {
		args[i] = a
	}
	return args
}

func (l *Logger) Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	// grab req_id from context
	args := grabRequestID(ctx, attrs)
	l.base.InfoContext(ctx, msg, args...)
}

func (l *Logger) Error(ctx context.Context, msg string, attrs ...slog.Attr) {
	args := grabRequestID(ctx, attrs)
	l.base.ErrorContext(ctx, msg, args...)
}
