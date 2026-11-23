package cli

import (
	"ProxyX/internal/platform"
	"fmt"
	"github.com/spf13/cobra"
	"os"
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
- 📄 Manage proxy configurations as YAML — list, apply, delete, and print samples
- 🌐 Edit configs in your browser — a localhost-only panel runs with the service
- 🔐 Issue and manage SSL/TLS certificates
- 📊 Monitor service health and status
- 🔄 Start, stop, and restart proxy instances

Common commands:
  proxyx config list          show configurations
  proxyx config apply <file>  apply a configuration
  proxyx config example       print a sample configuration
  proxyx start                start the proxy (serves the config panel at http://127.0.0.1:8787)

Designed for developers and sysadmins who need more than just a reverse proxy, ProxyX combines performance, flexibility, and ease-of-use in one tool.`,
	}

	cli := &CLI{
		Service:       service,
		root:          rootCmd,
		serviceConfig: "/etc/proxyx/conf.d",
	}

	// apply and delete are subcommands of `config` (see configCmd).
	cli.root.AddCommand(cli.configCmd())
	cli.root.AddCommand(cli.certCmd())
	cli.root.AddCommand(cli.healthCheckCmd())
	cli.root.AddCommand(cli.startCmd())
	cli.root.AddCommand(cli.restartCmd())
	cli.root.AddCommand(cli.stopCmd())
	cli.root.AddCommand(cli.statusCmd())
	cli.root.AddCommand(cli.updateCmd())
	cli.root.AddCommand(cli.versionCmd())
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

func requireRoot(cmd *cobra.Command, args []string) {
	if !isElevated() {
		fmt.Fprintln(os.Stderr, elevationMessage)
		os.Exit(1)
	}
}
