package cli

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)



func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}




var rootCmd = &cobra.Command{
	Use:   "proxyx",
	Short: "ProxyX CLI too and server",
}