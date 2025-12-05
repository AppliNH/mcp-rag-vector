package mcphandlers

import (
	"context"

	knowledgebase "github.com/applinh/mcp-rag-vector/gen/knowledge_base"
	knowledgebaseapi "github.com/applinh/mcp-rag-vector/internal/app/knowledgebase"
	"github.com/mark3labs/mcp-go/mcp"
)

func MCPKnowledgeBaseHandler(service knowledgebaseapi.Service) func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collection, err := req.RequireString("collection")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		content, err := req.RequireString("content")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var metadata map[string]any
		if args, ok := req.Params.Arguments.(map[string]any); ok {
			if metaArg, ok := args["metadata"]; ok {
				if m, ok := metaArg.(map[string]any); ok {
					metadata = m
				}
			}
		}

		success, err := service.Upsert(ctx, &knowledgebase.UpsertPayload{
			Collection: collection,
			Content:    content,
			Metadata:   metadata,
		})
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if success {
			return mcp.NewToolResultText("Successfully upserted document"), nil
		}
		return mcp.NewToolResultText("Upsert failed (unknown reason)"), nil
	}
}
