package cli

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/spf13/cobra"
)


func init() {
	rootCmd.AddCommand(reloadCmd)
}

func reloadProxyX() {
	cmd := exec.Command("sudo", "systemctl", "restart", "proxyx")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to reload ProxyX:", err)
		fmt.Println(string(output))
		return
	}

	fmt.Println("ProxyX reloaded successfully")
}

var reloadCmd = &cobra.Command{
	Use: "reload",
	Short: "Reload ProxyX configuration",
	Run: func(cmd *cobra.Command, args []string) {
		pidFile := "/var/run/proxyx.pid"
		pidBytes, err := os.ReadFile(pidFile)
		if err != nil {
			fmt.Println("ProxyX is not runnnig")
			return
		}

		var pid int
		fmt.Sscanf(string(pidBytes), "%d", &pid)
		process, _:= os.FindProcess(pid)
		err = process.Signal(syscall.SIGHUP)
		if err != nil {
			fmt.Println("Failed to reload ProxyX:", err)
		}

		fmt.Println("Reload Signal sent to ProxyX (SIGHUP)")
	},
}