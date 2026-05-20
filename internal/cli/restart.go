package cli

import (
	"github.com/spf13/cobra"
)

func (c *CLI) restartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart",
		PreRun: requireRoot,
		Short: "🔄 Restart the ProxyX service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Service.Restart()
		},
	}
}
