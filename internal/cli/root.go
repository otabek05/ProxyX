package cli

import (
	"ProxyX/internal/platform"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type CLI struct {
	Service platform.Service
	root    *cobra.Command
	serviceConfig string
	proxyXConfig string  

}

func NewCLI(service platform.Service) *CLI {

	var rootCmd = &cobra.Command{
		Use:   "proxyx",
		Short: "ProxyX CLI too and server",
	}

	cli := &CLI{
		Service: service,
		root: rootCmd,
		serviceConfig:  "/etc/proxyx/conf.d",

	}

	cli.root.AddCommand(cli.applyConfigCmd())
	cli.root.AddCommand(cli.deleteCmd())
	cli.root.AddCommand(cli.certCmd())

	cli.root.AddCommand(cli.configCmd())
	cli.root.AddCommand(cli.healthCheckCmd())
	cli.root.AddCommand(cli.restartCmd())
	cli.root.AddCommand(cli.stopCmd())
	cli.root.AddCommand(cli.statusCmd())


	return cli

}

func (c *CLI) Execute() {
	if err := c.root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
