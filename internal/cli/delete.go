package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete current configuration file",
	Example: `
     sudo proxyx delete local-proxy
     sudo proxyx delete my-api
  `,
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		configDir := "/etc/proxyx/conf.d"

		files, err := os.ReadDir(configDir)
		if err != nil {
			return fmt.Errorf("failed to read config directory: %v", err)
		}

		var matchedFile string

		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml") {

				fullPath := filepath.Join(configDir, file.Name())
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
			return fmt.Errorf("no configuration found with name: %s", name)
		}

		if err := os.Remove(matchedFile); err != nil {
			return fmt.Errorf("failed to delete: %v", err)
		}

		fmt.Printf("Deleted configuration '%s' (file: %s)\n", name, filepath.Base(matchedFile))
		restartProxyX()
		return nil
	},
}
