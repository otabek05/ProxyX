package cli

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)


func init() {
	rootCmd.AddCommand(reloadCmd)
}

func reloadProxyX() {
	cmd := exec.Command("sudo", "systemctl", "restart", "proxyx")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to restart ProxyX:", err)
		fmt.Println(string(output))
		return
	}

	fmt.Println("ProxyX restarted successfully")
}

var reloadCmd = &cobra.Command{
	Use: "restart",
	Short: "Reload ProxyX configuration",
	Run: func(cmd *cobra.Command, args []string) {
		reloadProxyX()
	},
}

