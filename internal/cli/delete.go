package cli

import (
	"ProxyX/internal/configmgr"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func (c *CLI) deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "delete [name]",
		Short:  "🗑️  Remove a ProxyX configuration",
		Args:   cobra.ExactArgs(1),
		PreRun: requireRoot,
		Example: `  sudo proxyx config delete local-proxy
  sudo proxyx config delete my-api`,
		Run: func(cmd *cobra.Command, args []string) {
			c.runDeleteFile(args[0])
		},
	}
}

func (c *CLI) runDeleteFile(name string) {
	file, err := configmgr.Delete(c.serviceConfig, name)
	if err != nil {
		color.Red.Printf("❌ %s\n", err.Error())
		os.Exit(1)
	}

	color.Green.Printf("🗑️  Configuration '%s' deleted (%s)\n", name, filepath.Base(file))
	fmt.Println("🔄 Restarting ProxyX service...")

	if err := c.Service.Restart(); err != nil {
		color.Red.Println("❌ Failed to restart ProxyX:", err)
		os.Exit(1)
	}

	color.Green.Println("✅ ProxyX restarted successfully")
}
