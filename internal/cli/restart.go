package cli

import (
	"github.com/spf13/cobra"
)

func (c *CLI) restartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		Short: "🔄 Restart the ProxyX service",
		Run: func(cmd *cobra.Command, args []string) {
			c.Service.Restart()
		},
	}
}
