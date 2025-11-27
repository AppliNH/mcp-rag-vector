package cmd

import (
	stdhttp "net/http"
	"sync"

	"github.com/applinh/mcp-rag-vector/cmd/services"
	infrahttp "github.com/applinh/mcp-rag-vector/internal/infra/http"

	"github.com/mark3labs/mcp-go/server"

	"github.com/spf13/cobra"

	goahttp "goa.design/goa/v3/http"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts API",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup
		ctx := cmd.Context()
		errc := make(chan error)

		mux := goahttp.NewMuxer()
		services.MountHealthService(ctx, mux, cfg)

		mcpSrv := mcpInit(ctx, cfg)

		// Use the MCP HTTP handler from mcp-go
		// Mount the Streamable HTTP server under /mcp
		mcpHandler := server.NewStreamableHTTPServer(mcpSrv, server.WithEndpointPath("/mcp"))

		// Mount MCP handler under /mcp and combine with existing goa mux
		finalHandler := infrahttp.CombineHandlers(mcpHandler, mux)

		// mount main handler under a small ServeMux
		top := stdhttp.NewServeMux()
		top.Handle("/", finalHandler)

		infrahttp.ServeHTTP(top, ctx, cfg, &wg, errc)
		for {
			if err := <-errc; err != nil {
				panic(err)
			}
			wg.Wait()
		}

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(mcpCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
