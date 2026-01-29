package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func (c *CLI) deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [name]",
		Short: "🗑️ Remove a ProxyX configuration or instance",
		PreRun: requireRoot,
		Example: `
     sudo proxyx delete local-proxy
     sudo proxyx delete my-api
  `,
		Run: func(cmd *cobra.Command, args []string)  {
			 c.runDeleteFile(args[0])
		},
	}
}
func (c *CLI) runDeleteFile(name string) {
	files, err := os.ReadDir(c.serviceConfig)
	if err != nil {
		fmt.Println("❌ Failed to read config directory:", err)
		return
	}

	var matchedFile string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {
			fullPath := filepath.Join(c.serviceConfig, file.Name())

			content, err := os.ReadFile(fullPath)
			if err != nil {
				continue
			}

			if strings.Contains(string(content), "name: "+name) {
				matchedFile = fullPath
				break
			}
		}
	}

	if matchedFile == "" {
		fmt.Printf("❌ Configuration '%s' not found\n", name)
		return
	}

	if err := os.Remove(matchedFile); err != nil {
		fmt.Println("❌ Failed to delete configuration:", err)
		return
	}

	fmt.Printf("🗑️ Configuration '%s' deleted (%s)\n", name, filepath.Base(matchedFile))
	fmt.Println("🔄 Restarting ProxyX service...")

	if err := c.Service.Restart(); err != nil {
		fmt.Println("❌ Failed to restart ProxyX:", err)
		return
	}

	fmt.Println("✅ ProxyX restarted successfully")
}
