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
	rootCmd.AddCommand(servicesCmd)
}

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Prints configured services by file",
	Run: func(cmd *cobra.Command, args []string) {
		configDir := "/etc/proxyx/configs"
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

				for _, route := range server.Spec.Routes {

					target := ""
					switch route.Type {
					case common.RouteReverseProxy:
						for _, url :=  range route.ReverseProxy.Servers {
							target += " , " + url.URL 
						}
					case common.RouteStatic:
						target = route.Static.Root
					}

					fmt.Printf(
						"%-20s %-25s %-10s %-40s\n",
						server.Spec.Domain,
						route.Path,
						route.Type,
						target,
					)
				}
		}
	},
}