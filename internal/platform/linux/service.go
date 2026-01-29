package linux

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Start() error {
	return nil
}

func (s *Service) Stop() error {
	cmd := exec.Command("sudo", "systemctl", "stop", "proxyx")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to stop ProxyX:", err)
		fmt.Println(string(output))
		return err
	}

	return nil
}

func (s *Service) Restart() error {
	cmd := exec.Command("sudo", "systemctl", "restart", "proxyx")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error: %v, Output: %s", err, string(output))
	}

	return nil
}

func (s *Service) Status() error {
	cmdCheck := exec.Command("systemctl", "is-active", "proxyx")
	output, _ := cmdCheck.CombinedOutput()
	status := strings.TrimSpace(string(output))

	if status != "active" {
		fmt.Println("ProxyX is not running")
		return fmt.Errorf("ProxyX is not running")
	}

	fmt.Println("ProxyX is running (systemd service)")

	cmdPID := exec.Command("systemctl", "show", "proxyx", "-p", "MainPID")
	pidOutput, _ := cmdPID.CombinedOutput()
	var pid int
	fmt.Sscanf(string(pidOutput), "MainPID=%d", &pid)
	if pid == 0 {
		return fmt.Errorf("Cannot find ProxyX PID")
	}

	psCmd := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "pid,pcpu,pmem,etime=", "--no-headers")
	psOutput, _ := psCmd.CombinedOutput()
	psFields := strings.Fields(string(psOutput))
	if len(psFields) < 4 {
		return fmt.Errorf("Failed to get process stats")
	}


	fmt.Println("📊 ProxyX Process Status")
	fmt.Printf("    PID     : %s\n", psFields[0])
	fmt.Printf("    CPU     : %s %%\n", psFields[1])
	fmt.Printf("    Memory  : %s %%\n", psFields[2])
	fmt.Printf("    Uptime  : %s\n", psFields[3])
	return nil
}
