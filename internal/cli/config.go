package cli

import (
	"ProxyX/internal/common"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	//"github.com/olekukonko/tablewriter"
	//"github.com/olekukonko/tablewriter/renderer"
	//"github.com/olekukonko/tablewriter/tw"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func (c *CLI) configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "⚙️ Display current proxy configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			return c.runConfigs(output)
		},
	}

	cmd.Flags().StringP("output", "o", "", "Output format (wide|default)")
	return cmd

}

func (c *CLI) runConfigs(output string) error {
	files, err := filepath.Glob(filepath.Join(c.serviceConfig, "*.yaml"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		fmt.Println("ℹ️ No configuration files available.")
		return nil
	}

	wide := output == "wide"

	fmt.Println("📄 Loaded Proxy Configurations")
	fmt.Println(strings.Repeat("=", 40))

	configCount := 0

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			color.Red.Println("Failed to read file:", file)
			continue
		}

		var server common.ServerConfig
		if err := yaml.Unmarshal(data, &server); err != nil {
			color.Red.Println("Invalid YAML:", file)
			continue
		}

		configCount++

		// ── Basic info (always shown)
		fmt.Printf("\n📦 Config Name : %s\n", server.Metadata.Name)
		fmt.Printf("🌍 Domain      : %s\n", server.Spec.Domain)

		// ── Wide-only info
		if wide {
			fmt.Printf("🗂️  Namespace  : %s\n", server.Metadata.Namespace)
			fmt.Printf("📁 File        : %s\n", filepath.Base(file))

			// TLS
			if server.Spec.TLS != nil {
				fmt.Println("🔐 TLS")
				fmt.Printf("    Cert : %s\n", server.Spec.TLS.CertFile)
				fmt.Printf("    Key  : %s\n", server.Spec.TLS.KeyFile)
			}

			// Rate limit
			rl := server.Spec.RateLimit
			if rl.Requests > 0 {
				fmt.Println("🚦 Rate Limit")
				fmt.Printf("    %d req / %ds\n", rl.Requests, rl.WindowSeconds)
			}
		}

	//	fmt.Println()

		// ── Routes (always shown)
		for _, route := range server.Spec.Routes {
			fmt.Println("📍 Route")
			fmt.Printf("    Path : %s\n", route.Path)
			fmt.Printf("    Type : %s\n", route.Type.String())

			switch route.Type {
			case common.RouteReverseProxy:
				fmt.Println("    Target :")
				for _, s := range route.ReverseProxy.Servers {
					fmt.Printf("        → %s\n", s.URL)
				}

			case common.RouteStatic:
				fmt.Printf("    Target : %s\n", route.Static.Root)

			case common.RouteWebsocket:
				fmt.Printf("    Target : %s\n", route.Websocket.URL)
			}

//			fmt.Println()
		}

		fmt.Println(strings.Repeat("-", 40))
	}

	fmt.Printf("✨ Total Configs : %d\n", configCount)
	fmt.Println(strings.Repeat("=", 40))

	return nil
}




/*

func (c *CLI) runConfigs(output string) error {
		files, err := filepath.Glob(filepath.Join(c.serviceConfig, "*.yaml"))
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Println("ℹ️ No configuration files available.")
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
				case common.RouteWebsocket:
					target = route.Websocket.URL
				}

				if output == "wide" {

					rl := server.Spec.RateLimit
					var tls string 
					if server.Spec.TLS != nil {
						tls = fmt.Sprintf("%s\n%s", server.Spec.TLS.CertFile, server.Spec.TLS.KeyFile)
					}

					table.Append([]string{
						filepath.Base(file),
						server.Metadata.Name,
						server.Metadata.Namespace,
						server.Spec.Domain,
						route.Path,
						route.Type.String(),
						target,
						fmt.Sprintf("%d req / %ds", rl.Requests, rl.WindowSeconds),
						tls,
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
} 


*/
