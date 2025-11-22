package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/applinh/mcp-rag-vector/cmd/config"
	"github.com/applinh/mcp-rag-vector/cmd/services"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Starts the app in MCP mode",
	Run: func(cmd *cobra.Command, args []string) {

		ctx := cmd.Context()
		mcpSrv := mcpInit(ctx, cfg)

		// If requested, run the MCP server over stdio (useful when Claude
		// Desktop spawns the process and expects stdio-based MCP transport).

		if err := server.ServeStdio(mcpSrv); err != nil {
			fmt.Fprintf(os.Stderr, "MCP stdio server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func mcpInit(ctx context.Context, cfg config.Config) *server.MCPServer {
	mcpSrv := server.NewMCPServer("MCP", "1.0.0")
	services.MountGreetingMCPService(ctx, mcpSrv, cfg)

	return mcpSrv
}
