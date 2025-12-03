package cli

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stopCmd)
}

var stopCmd = &cobra.Command{
	Use: "stop",
	Short: "Stops running proxyx service",
	Run: func(cmd *cobra.Command, args []string) {
		restartCMD := exec.Command("sudo", "systemctl", "stop", "proxyx")
		output , err := restartCMD.CombinedOutput()
		if err != nil {
		fmt.Println("Failed to restart ProxyX:", err)
		fmt.Println(string(output))
		return
	}

	fmt.Println("ProxyX stopped successfully")
	},
}