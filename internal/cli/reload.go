package cli

import (
	"github.com/spf13/cobra"
)

func (c *CLI) restartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "Reload ProxyX configuration",
		Run: func(cmd *cobra.Command, args []string) {
			c.Service.Restart()
		},
	}
}
