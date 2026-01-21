package cli

import (
	"ProxyX/internal/proxy"
	"ProxyX/pkg/config"
	"log"

	"github.com/spf13/cobra"
)

func (c *CLI) runCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Runs proxyx in the background if not running",
		RunE: func(cmd *cobra.Command, args []string) error {
			
			proxyConfig, err := config.LoadProxyXConfig()
			if err != nil {
				log.Fatalf("Failed to load proxy config: %v", err)
			}

			serverConfig, err := config.LoadConfig()
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}

			srv := proxy.NewServer(serverConfig, proxyConfig)
			return srv.Start()
		},
	}

}
