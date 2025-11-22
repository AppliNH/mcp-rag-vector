package mcphandlers

import (
	"context"

	"github.com/applinh/mcp-rag-vector/gen/greeting"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCPGreetingHandler adapts the Goa greeting service for MCP.
func MCPGreetingHandler(service greeting.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := req.RequireString("name")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		greetingMsg, err := service.Greet(ctx, &greeting.GreetPayload{
			Name: name,
		})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		return mcp.NewToolResultText(greetingMsg), nil
	}
}
