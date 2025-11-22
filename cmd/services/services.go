package services

import (
	"context"

	"github.com/applinh/mcp-rag-vector/cmd/config"
	"github.com/applinh/mcp-rag-vector/gen/health"
	greetingapi "github.com/applinh/mcp-rag-vector/internal/app/greeting"
	healthapi "github.com/applinh/mcp-rag-vector/internal/app/health"
	"github.com/applinh/mcp-rag-vector/internal/infra/http"
	"github.com/applinh/mcp-rag-vector/internal/infra/logger"
	mcphandlers "github.com/applinh/mcp-rag-vector/internal/infra/mcp_handlers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

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
	loggerInstance := logger.NewLogger(slogger)
	healthSvc := healthapi.NewHealthService(loggerInstance)
	heathEndpoints := health.NewEndpoints(healthSvc)
	handler := healthsvr.New(heathEndpoints, mux, goahttp.RequestDecoder, goahttp.ResponseEncoder, http.ErrorHandler(loggerInstance), defaultFormat)

	healthsvr.Mount(mux, handler)
}

func MountGreetingMCPService(ctx context.Context, mcpSrv *server.MCPServer, cfg config.Config) {
	slogger := defaultSLoggerSettings("greeting", cfg.LogLevel)
	loggerInstance := logger.NewLogger(slogger)
	greetingSvc := greetingapi.NewGreetingService(loggerInstance)
	greetingTool := mcp.NewTool("greet",
		mcp.WithDescription("Greeting tool"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Name to greet")),
	)
	mcpSrv.AddTool(greetingTool, mcphandlers.MCPGreetingHandler(greetingSvc))
}
