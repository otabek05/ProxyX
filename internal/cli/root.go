package cli

import (
	"ProxyX/internal/platform"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type CLI struct {
	Service       platform.Service
	root          *cobra.Command
	serviceConfig string
	proxyXConfig  string
}

func NewCLI(service platform.Service) *CLI {

	var rootCmd = &cobra.Command{
		Use:   "proxyx",
		Short: "🚀 ProxyX — Advanced proxy management CLI",
		Long: `🚀 ProxyX is a powerful and flexible command-line tool for managing proxies, similar to Nginx but with extended capabilities.

With ProxyX, you can:
- Effortlessly manage proxy configurations and SSL/TLS certificates
- Monitor service health and status
- Start, stop, and restart proxy instances
- Gain full control over your proxy environment with a simple CLI

Designed for developers and sysadmins who need more than just a reverse proxy, ProxyX combines performance, flexibility, and ease-of-use in one tool.`,
	}

	cli := &CLI{
		Service:       service,
		root:          rootCmd,
		serviceConfig: "/etc/proxyx/conf.d",
	}

	cli.root.AddCommand(cli.applyConfigCmd())
	cli.root.AddCommand(cli.deleteCmd())
	cli.root.AddCommand(cli.certCmd())

	cli.root.AddCommand(cli.configCmd())
	cli.root.AddCommand(cli.healthCheckCmd())
	cli.root.AddCommand(cli.restartCmd())
	cli.root.AddCommand(cli.stopCmd())
	cli.root.AddCommand(cli.statusCmd())
	cli.root.AddCommand(cli.versionCmd())
	cli.root.AddCommand(cli.startCmd())

	return cli

}

func (c *CLI) Execute() {
	if err := c.root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (c *CLI) Help() {
	c.root.Help()
}
