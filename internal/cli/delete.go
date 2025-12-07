package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.Flags().StringP("file", "f", "", "Config file to delete")
	deleteCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete current configuration file",
	Example: `
     sudo proxyx delete -f local-proxy.yaml
     sudo proxyx delete --file test.yaml
  `,
	RunE: func(cmd *cobra.Command, args []string) error {
		fileName, _ := cmd.Flags().GetString("file")
		configDir := "/etc/proxyx/configs"

		if strings.Contains(fileName, "..") || strings.Contains(fileName, "/") {
			return fmt.Errorf("invalid file name: %s ", fileName)
		}

		fullPath := filepath.Join(configDir, fileName)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("config file does not exist: %s", fileName)
		} 

		if err := os.Remove(fullPath); err != nil {
			return fmt.Errorf("failed to delete file: %s", fileName)
		}

		fmt.Printf("Deleted config: %s\n", fileName)
		return nil
	},
}
