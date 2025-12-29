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
		Short: "Delete current configuration file",
		Example: `
     sudo proxyx delete local-proxy
     sudo proxyx delete my-api
  `,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.runDeleteFile(args[0])
		},
	}
}

func (c *CLI) runDeleteFile(name string) error {
			files, err := os.ReadDir(c.serviceConfig)
			if err != nil {
				return fmt.Errorf("failed to read config directory: %v", err)
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
				return fmt.Errorf("no configuration found with name: %s", name)
			}

			if err := os.Remove(matchedFile); err != nil {
				return fmt.Errorf("failed to delete: %v", err)
			}

			fmt.Printf("Deleted configuration '%s' (file: %s)\n", name, filepath.Base(matchedFile))
			return c.Service.Restart()
		
}
