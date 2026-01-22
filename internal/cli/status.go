package cli

import (
	"github.com/spf13/cobra"
)

func (c *CLI) statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "📊 Display ProxyX runtime status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Service.Status()
		},
	}
}