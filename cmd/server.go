package cmd

import (
	"sync"

	"github.com/applinh/mcp-rag-vector/cmd/services"
	"github.com/applinh/mcp-rag-vector/internal/infra/http"

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

		http.ServeHTTP(mux, ctx, cfg, &wg, errc)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
