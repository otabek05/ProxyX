package cli

import (
	"ProxyX/pkg/config"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	hcEnable   bool
	hcDisable  bool
	hcPath     string
	hcInterval time.Duration
)

func (c *CLI) healthCheckCmd() *cobra.Command {
	healthCheckCmd := &cobra.Command{
		Use:   "healthcheck",
		Short: "Configure health check settings",
		Long: `Enable, disable, or configure the health check endpoint.

Examples:
  proxyx healthcheck --enable
  proxyx healthcheck --disable
  proxyx healthcheck --path /health
  proxyx healthcheck --interval 10s`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runHealthCheck(cmd)
		},
	}

	healthCheckCmd.Flags().BoolVar(&hcEnable, "enable", false, "Enable health check")
	healthCheckCmd.Flags().BoolVar(&hcDisable, "disable", false, "Disable health check")
	healthCheckCmd.Flags().StringVar(&hcPath, "path", "/health", "Health check path")
	healthCheckCmd.Flags().DurationVar(&hcInterval, "interval", 10*time.Second, "Health check interval")

	return healthCheckCmd
}

func (c *CLI) runHealthCheck(cmd *cobra.Command) error {
	cfg, err := config.LoadProxyXConfig()
	if err != nil {
		return err
	}

	updated := false

	// Enable / Disable
	if hcEnable && hcDisable {
		return fmt.Errorf("cannot use --enable and --disable together")
	}

	if hcEnable {
		cfg.HealthCheck.Enabled = true
		fmt.Println("✔ Health check enabled")
		updated = true
	}

	if hcDisable {
		cfg.HealthCheck.Enabled = false
		fmt.Println("✔ Health check disabled")
		updated = true
	}

	// Path
	if cmd.Flags().Changed("path") {
		cfg.HealthCheck.Path = hcPath
		fmt.Printf("✔ Health check path set to %s\n", hcPath)
		updated = true
	}

	// Interval
	if cmd.Flags().Changed("interval") {
		cfg.HealthCheck.Interval = hcInterval
		fmt.Printf("✔ Health check interval set to %s\n", hcInterval)
		updated = true
	}

	if !updated {
		fmt.Println("ℹ No changes applied")
		fmt.Printf("Current settings:\n")
		fmt.Printf("  Enabled : %v\n", cfg.HealthCheck.Enabled)
		fmt.Printf("  Path    : %s\n", cfg.HealthCheck.Path)
		fmt.Printf("  Interval: %s\n", cfg.HealthCheck.Interval)
		return nil
	}

	if err := config.SaveProxyXConfig(cfg); err != nil {
		return err
	}

	return c.Service.Restart()

}
