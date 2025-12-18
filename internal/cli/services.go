package cli

import (
	"ProxyX/internal/common"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func init() {
	configCmd.Flags().StringP("output", "o", "", "Output format")
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "configs",
	Short: "Print configured proxy configs",
	RunE: func(cmd *cobra.Command, args []string) error {
		output, _ := cmd.Flags().GetString("output")
		configDir := "/etc/proxyx/conf.d"

		files, err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Println("No configuration files found.")
			return nil
		}

		table := tablewriter.NewTable(os.Stdout, tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Settings: tw.Settings{Separators: tw.Separators{BetweenRows: tw.On}},
		})),
			tablewriter.WithConfig(tablewriter.Config{
				Header: tw.CellConfig{Alignment: tw.CellAlignment{Global: tw.AlignCenter}},
				Row: tw.CellConfig{
					Merging:   tw.CellMerging{Mode: tw.MergeHierarchical},
					Alignment: tw.CellAlignment{Global: tw.AlignLeft},
				},
			}),
		)

		if output == "wide" {
			table.Header([]string{
				"FILE", "NAME", "NAMESPACE",
				"DOMAIN", "PATH", "TYPE", "TARGET",
				"RATELIMIT", "TLS",
			})
		} else {
			table.Header([]string{"NAME", "DOMAIN", "PATH", "TYPE", "TARGET"})
		}

		for _, file := range files {

			data, _ := os.ReadFile(file)
			var server common.ServerConfig
			if err := yaml.Unmarshal(data, &server); err != nil {
				color.Red.Println("Invalid YAML:", file)
				continue
			}

			for _, route := range server.Spec.Routes {

				var target string
				switch route.Type {
				case common.RouteReverseProxy:
					if len(route.ReverseProxy.Servers) == 1 {
						target = route.ReverseProxy.Servers[0].URL
					} else {
						lines := []string{}
						for _, s := range route.ReverseProxy.Servers {
							lines = append(lines, s.URL)
						}
						target = strings.Join(lines, "\n")
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
						fmt.Sprintf("%s\n%s", tls.CertFile, tls.KeyFile),
					})

				} else {
					table.Append([]string{
						server.Metadata.Name,
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
