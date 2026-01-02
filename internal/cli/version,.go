package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c *CLI) versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show ProxyX version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ProxyX version v0.1.3")
		},
	}
}

