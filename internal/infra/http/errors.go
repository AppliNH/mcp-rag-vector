package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/applinh/mcp-rag-vector/internal/infra/logger"
)

func ErrorHandler(log logger.LoggerInterface) func(context.Context, http.ResponseWriter, error) {
	return func(ctx context.Context, w http.ResponseWriter, err error) {
		log.Error(ctx, "internal server error occurred", slog.String("error", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)

		// write json response
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"error": "internal server error"}`)); err != nil {
			log.Error(ctx, "fatal: failed to write error response", slog.String("error", err.Error()))
		}

	}
}
