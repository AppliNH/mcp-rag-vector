package services

import (
	"log/slog"
	"os"
)

func defaultSLoggerSettings(serviceName string, logLevel slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	})

	return slog.New(handler).With("service", serviceName)
}
