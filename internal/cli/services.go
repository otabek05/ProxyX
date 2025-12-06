package cli

import (
	"ProxyX/internal/common"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)


func init() {
	servicesCmd.Flags().StringP("delete", "d", "", "Delete a config file")
	rootCmd.AddCommand(servicesCmd)
}

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Prints configured services by file",
	Run: func(cmd *cobra.Command, args []string) {

		deleteFile, _ := cmd.Flags().GetString("delete")
		configDir := "/etc/proxyx/configs"
		if deleteFile != "" {
			fullPath := filepath.Join(configDir, deleteFile)

			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				fmt.Println("Config file does not exist:", deleteFile)
		        return
			}else if err != nil {
				fmt.Println("Failed to check file: ", err)
				return
			}

			if err := os.Remove(fullPath); err != nil {
				fmt.Println("Failed to delete file:", err)
				return
			}

			fmt.Printf("Deleted config: %s\n", deleteFile)
			return
		}

		files, err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
		if err != nil {
			fmt.Println("Failed to list configs:", err)
			return
		}

		if len(files) == 0 {
			fmt.Println("No configuration files found.")
			return
		}

		for _, file := range files {
			fmt.Println("\n📄 File:", filepath.Base(file))
			fmt.Println(strings.Repeat("=", 90))

			data, err := os.ReadFile(file)
			if err != nil {
				fmt.Println("Failed to read:", file)
				continue
			}

			var server common.ServerConfig
			if err := yaml.Unmarshal(data, &server); err != nil {
				fmt.Println("Invalid YAML:", file)
				continue
			}

			fmt.Printf("%-20s %-25s %-10s %-40s\n", "DOMAIN", "PATH", "TYPE", "TARGET")
			fmt.Println(strings.Repeat("-", 95))

				for _, route := range server.Routes {

					target := ""
					switch route.Type {
					case "proxy":
						target = strings.Join(route.Backends, ", ")
					case "static":
						target = route.Dir
					default:
						target = "unknown"
					}

					fmt.Printf(
						"%-20s %-25s %-10s %-40s\n",
						server.Domain,
						route.Path,
						route.Type,
						target,
					)
				}
		}
	},
}