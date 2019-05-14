package apiserver

import (
	"github.com/nsyszr/ariane/pkg/cmd/apiserver/server"
	"github.com/spf13/cobra"
)

// servePublicCmd represents the serve public command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run API Server",
	Run:   server.RunServe(c),
}

func init() {
	RootCmd.AddCommand(serveCmd)

	/*serveAllCmd.Flags().BoolVarP(&c.SkipTLSVerify, "skip-tls-verify", "", false, "Skip TLS verification on HTTP client requests.")
	serveAllCmd.Flags().BoolVarP(&c.ForceHTTP, "dangerous-force-http", "", false, "Do not run this in production")
	serveAllCmd.Flags().BoolVarP(&c.DevMode, "dev-mode", "", false, "Do not run this in production")*/
}
