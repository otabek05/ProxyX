package cli

import (
	"ProxyX/internal/common"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
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
	files, err := os.ReadDir(c.serviceConfig)
	if err != nil {
		color.Red.Println("❌ Failed to read config directory:", err)
		os.Exit(1)
	}

	var matchedFile string
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".yaml") && !strings.HasSuffix(file.Name(), ".yml") {
			continue
		}
		fullPath := filepath.Join(c.serviceConfig, file.Name())

		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		var server common.ServerConfig
		if err := yaml.Unmarshal(content, &server); err != nil {
			continue
		}
		if server.Metadata.Name == name {
			matchedFile = fullPath
			break
		}
	}

	if matchedFile == "" {
		color.Red.Printf("❌ Configuration '%s' not found\n", name)
		os.Exit(1)
	}

	if err := os.Remove(matchedFile); err != nil {
		color.Red.Println("❌ Failed to delete configuration:", err)
		os.Exit(1)
	}

	color.Green.Printf("🗑️  Configuration '%s' deleted (%s)\n", name, filepath.Base(matchedFile))
	fmt.Println("🔄 Restarting ProxyX service...")

	if err := c.Service.Restart(); err != nil {
		color.Red.Println("❌ Failed to restart ProxyX:", err)
		os.Exit(1)
	}

	color.Green.Println("✅ ProxyX restarted successfully")
}
