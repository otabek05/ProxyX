package cli

import (
	"github.com/spf13/cobra"
)

func (c *CLI) stopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stops running proxyx service",
		RunE: func(cmd *cobra.Command, args []string) error  {
			return c.Service.Stop()
		},
	}

}

