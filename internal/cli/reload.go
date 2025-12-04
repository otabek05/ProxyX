package cli

import (
	"fmt"
	"os/exec"
	"github.com/spf13/cobra"
)


func init(){
	rootCmd.AddCommand(restartCmd)
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

var restartCmd = &cobra.Command{
	Use: "restart",
	Short: "Reload ProxyX configuration",
	Run: func(cmd *cobra.Command, args []string) {
		reloadProxyX()
	},
}

