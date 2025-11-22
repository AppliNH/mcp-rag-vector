package http

import (
	"net/http"
	"strings"
)

// CombineHandlers returns an http.Handler that dispatches requests whose
// path starts with "/mcp" to the provided mcpHandler. All other requests
// are forwarded to mainHandler.
//
// If mcpHandler is nil the returned handler simply forwards all requests to
// mainHandler.
func CombineHandlers(mcpHandler http.Handler, mainHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mcpHandler != nil && strings.HasPrefix(r.URL.Path, "/mcp") {
			// If the MCP handler expects requests at the root path, you may
			// want to use http.StripPrefix("/mcp", mcpHandler) when creating
			// the handler. Keep behaviour explicit here: forward as-is.
			mcpHandler.ServeHTTP(w, r)
			return
		}
		mainHandler.ServeHTTP(w, r)
	})
}
