package cli

import (
	"ProxyX/internal/proxy"
	"ProxyX/internal/pkg/config"
	"log"
	"net"

	"github.com/spf13/cobra"
)

func (c *CLI) startCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "start",
		Short:  "▶️ Start the ProxyX service",
		PreRun: requireRoot,
		Run: func(cmd *cobra.Command, args []string) {

			if IsServerRunning() {
				log.Println("ProxyX server is already running.")
				return
			}

			proxyConfig, err := config.LoadProxyXConfig()
			if err != nil {
				log.Fatalf("Failed to load proxy config: %v", err)
			}

			serverConfig, err := config.LoadConfig()
			if err != nil {
				log.Fatalf("Failed to load config: %v", err)
			}

			srv := proxy.NewServer(serverConfig, proxyConfig)
			runServer(srv)
		},
	}

}

func IsServerRunning() bool {
	conn, err := net.Dial("tcp", "127.0.0.1:443")
	if err != nil {
		return false
	}

	conn.Close()
	return true
}
