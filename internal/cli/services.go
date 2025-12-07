package cli

import (
	"ProxyX/internal/common"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	servicesCmd.Flags().StringP("output", "o", "", "Outpuy")
	rootCmd.AddCommand(servicesCmd)
}

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "Print configured services by file",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		configDir := "/etc/proxyx/configs"
		files, err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Println("No configuration files found.")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)

		// Set headers depending on -o wide
		if output == "wide" {
			table.Header([]string{"FILE", "NAME", "NAMESPACE", "DOMAIN", "PATH", "TYPE", "TARGET", "RATELIMIT", "TLS"})
		} else {
			table.Header([]string{"DOMAIN", "PATH", "TYPE", "TARGET"})
		}

		for _, file := range files {
			data, _ := os.ReadFile(file)
			var server common.ServerConfig
			if err := yaml.Unmarshal(data, &server); err != nil {
				color.Red.Println("Invalid YAML:", file)
				continue
			}

			for _, route := range server.Spec.Routes {
				target := ""
				switch route.Type {
				case common.RouteReverseProxy:
					if len(route.ReverseProxy.Servers) == 1 {
						target = route.ReverseProxy.Servers[0].URL
					} else {
						var parts []string
						for _, s := range route.ReverseProxy.Servers {
							parts = append(parts, s.URL)
						}
						target = strings.Join(parts, "\n")
					}
				case common.RouteStatic:
					target = route.Static.Root
				}

				if output == "wide" {
					rl := server.Spec.RateLimit
					tls := server.Spec.TLS
					table.Append([]string{
						filepath.Base(file),
						server.Metadata.Name,
						server.Metadata.Namespace,
						server.Spec.Domain,
						route.Path,
						route.Type.String(),
						target,
						fmt.Sprintf("%d req / %ds", rl.Requests, rl.WindowSeconds),
						fmt.Sprintf("%s \n %s", tls.CertFile, tls.KeyFile),
					})
				} else {
					table.Append([]string{
						server.Spec.Domain,
						route.Path,
						route.Type.String(),
						target,
					})
				}
			}
		}

		table.Render()
		return nil
	},
}
