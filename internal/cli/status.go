package cli

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if ProxyX is running or not",
	Run: func(cmd *cobra.Command, args []string) {
		cmdCheck := exec.Command("systemctl", "is-active", "proxyx")
		output, _ := cmdCheck.CombinedOutput()
		status := strings.TrimSpace(string(output))

		if status != "active" {
			fmt.Println("ProxyX is not running")
			return
		}

		fmt.Println("ProxyX is running (systemd service)")

		cmdPID := exec.Command("systemctl", "show", "proxyx", "-p", "MainPID")
		pidOutput, _ := cmdPID.CombinedOutput()
		var pid int
		fmt.Sscanf(string(pidOutput), "MainPID=%d", &pid)
		if pid == 0 {
			fmt.Println("Cannot find ProxyX PID")
			return
		}

		psCmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "pid,pcpu,pmem,etime=", "--no-headers")
		psOutput, _ := psCmd.CombinedOutput()
		psFields := strings.Fields(string(psOutput))
		if len(psFields) < 4 {
			fmt.Println("Failed to get process stats")
			return
		}

		pidStrOut := psFields[0]
		cpu := psFields[1]
		mem := psFields[2]
		uptime := psFields[3]

		fmt.Println("PID       CPU%    MEM%    Uptime")
		fmt.Printf("%-9s %-7s %-7s %-7s\n", pidStrOut, cpu, mem, uptime)
	},
}
