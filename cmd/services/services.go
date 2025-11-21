package services

import (
	"context"

	"github.com/applinh/mcp-rag-vector/cmd/config"
	"github.com/applinh/mcp-rag-vector/gen/health"
	healthapi "github.com/applinh/mcp-rag-vector/internal/app/health"
	"github.com/applinh/mcp-rag-vector/internal/infra/http"
	"github.com/applinh/mcp-rag-vector/internal/infra/logger"

	sloghttp "github.com/samber/slog-http"

	healthsvr "github.com/applinh/mcp-rag-vector/gen/http/health/server"

	goahttp "goa.design/goa/v3/http"
)

// Add request ID in errors payload to match in logs - can later be used as TraceID
var defaultFormat = func(ctx context.Context, err error) goahttp.Statuser {
	se := goahttp.NewErrorResponse(ctx, err)
	if se, ok := se.(*goahttp.ErrorResponse); ok {
		se.ID = sloghttp.GetRequestIDFromContext(ctx)
	}
	return se
}

func MountHealthService(ctx context.Context, mux goahttp.Muxer, cfg config.Config) {
	slogger := defaultSLoggerSettings("greeting", cfg.LogLevel)
	logger := logger.NewLogger(slogger)
	healthSvc := healthapi.NewHealthService(logger)
	heathEndpoints := health.NewEndpoints(healthSvc)
	handler := healthsvr.New(heathEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, http.ErrorHandler(logger), defaultFormat)

	healthsvr.Mount(mux, handler)
}
