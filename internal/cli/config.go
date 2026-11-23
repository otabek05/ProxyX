package cli

import (
	"ProxyX/internal/common"
	"ProxyX/internal/configmgr"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// configCmd is the parent for all configuration management commands:
//
//	proxyx config list       show loaded configurations
//	proxyx config apply      apply a YAML configuration
//	proxyx config delete     remove a configuration
//	proxyx config example    print a sample configuration
//
// Running `proxyx config` with no subcommand lists configurations.
func (c *CLI) configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "⚙️  Manage ProxyX configurations",
		Long:  "⚙️  Manage ProxyX configurations: list, apply, delete, and print a sample.",
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			return c.runConfigs(output)
		},
	}
	cmd.Flags().StringP("output", "o", "", "Output format (wide|default)")

	cmd.AddCommand(c.configListCmd())
	cmd.AddCommand(c.applyConfigCmd())
	cmd.AddCommand(c.deleteCmd())
	cmd.AddCommand(c.configExampleCmd())
	return cmd
}

func (c *CLI) configListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "📄 List loaded proxy configurations",
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
		color.Yellow.Println("ℹ️  No configuration files available.")
		return nil
	}

	wide := output == "wide"

	color.Cyan.Println("📄 Loaded Proxy Configurations")
	fmt.Println(strings.Repeat("=", 40))

	configCount := 0

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			color.Red.Println("❌ Failed to read file:", file)
			continue
		}

		var server common.ServerConfig
		if err := yaml.Unmarshal(data, &server); err != nil {
			color.Red.Println("❌ Invalid YAML:", file)
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

			if server.Spec.TLS != nil {
				fmt.Println("🔐 TLS")
				fmt.Printf("    Cert : %s\n", server.Spec.TLS.CertFile)
				fmt.Printf("    Key  : %s\n", server.Spec.TLS.KeyFile)
			}

			if rl := server.Spec.RateLimit; rl != nil && rl.Requests > 0 {
				fmt.Println("🚦 Rate Limit")
				fmt.Printf("    %d req / %ds\n", rl.Requests, rl.WindowSeconds)
			} else {
				fmt.Println("🚦 Rate Limit : global default")
			}
		}

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
		}

		fmt.Println(strings.Repeat("-", 40))
	}

	color.Green.Printf("✨ Total Configs : %d\n", configCount)
	fmt.Println(strings.Repeat("=", 40))

	return nil
}

func (c *CLI) configExampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "example",
		Short: "📝 Print a sample configuration you can save and apply",
		Long: "📝 Print a sample configuration to stdout.\n\n" +
			"By default it prints a per-domain server config. Use --global to print\n" +
			"the global config.yaml (timeouts, default rate limit, IP whitelist).",
		Example: `  proxyx config example > my-app.yaml
  sudo proxyx config apply my-app.yaml

  proxyx config example --global`,
		RunE: func(cmd *cobra.Command, args []string) error {
			global, _ := cmd.Flags().GetBool("global")
			if global {
				fmt.Print(configmgr.GlobalExample)
			} else {
				fmt.Print(configmgr.ServerExample)
			}
			return nil
		},
	}
	cmd.Flags().Bool("global", false, "Print the global config.yaml example instead of a server config")
	return cmd
}
