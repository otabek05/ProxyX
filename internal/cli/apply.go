package cli

import (
	"ProxyX/internal/configmgr"
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func (c *CLI) applyConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "apply [file]",
		Short:   "⚡ Apply a YAML configuration to ProxyX",
		Args:    cobra.ExactArgs(1),
		PreRun:  requireRoot,
		Example: "  sudo proxyx config apply ./my-app.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.runApply(args[0]); err != nil {
				color.Red.Println(err.Error())
				os.Exit(1)
			}
		},
	}
}

func (c *CLI) runApply(file string) error {
	if file == "" {
		return fmt.Errorf("❌ please provide a config file path, e.g. proxyx config apply ./my-app.yaml")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("❌ cannot read file: %w", err)
	}

	if _, err := configmgr.Apply(c.serviceConfig, data); err != nil {
		return fmt.Errorf("❌ %w", err)
	}

	color.Green.Println("✅ Configuration applied successfully")
	fmt.Println("🔄 Restarting ProxyX service...")

	if err := c.Service.Restart(); err != nil {
		return fmt.Errorf("❌ failed to restart ProxyX: %w", err)
	}

	color.Green.Println("✅ ProxyX restarted successfully")
	return nil
}
