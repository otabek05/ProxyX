package cli

import "github.com/spf13/cobra"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Show ProxyX version",
}