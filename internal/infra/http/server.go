package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/applinh/mcp-rag-vector/cmd/config"

	sloghttp "github.com/samber/slog-http"

	goahttp "goa.design/goa/v3/http"
)

var (
	Decoder = goahttp.RequestDecoder
	Encoder = goahttp.ResponseEncoder
)

func ServeHTTP(mux goahttp.Muxer, ctx context.Context, cfg config.Config, wg *sync.WaitGroup, errc chan error) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.LogLevel,
	}))

	// Middleware
	handler := sloghttp.Recovery(mux)
	handler = sloghttp.New(logger)(handler)

	// Start HTTP server using default configuration, change the code to
	// configure the server as required by your service.
	srv := &http.Server{Addr: "0.0.0.0:" + cfg.HTTP.Port, Handler: handler, ReadHeaderTimeout: time.Second * 60}
	(*wg).Add(1)
	go func() {
		defer (*wg).Done()

		// Start HTTP server in a separate goroutine.
		go func() {
			logger.InfoContext(ctx, fmt.Sprintf("HTTP server listening on port %v", cfg.HTTP.Port))
			errc <- srv.ListenAndServe()
		}()

		<-ctx.Done()
		logger.InfoContext(ctx, fmt.Sprintf("shutting down HTTP server at %v", cfg.HTTP.Port))

		// Shutdown gracefully with a 30s timeout.
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			logger.ErrorContext(ctx, "failed to shutdown", slog.Any("error", err))
		}
	}()
}
