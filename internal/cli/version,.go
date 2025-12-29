package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)


var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Show ProxyX version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ProxyX version 1.0.0")
	},
}