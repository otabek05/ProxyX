package cli

import (
	"ProxyX/internal/common"
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
				fmt.Print(globalConfigExample)
			} else {
				fmt.Print(serverConfigExample)
			}
			return nil
		},
	}
	cmd.Flags().Bool("global", false, "Print the global config.yaml example instead of a server config")
	return cmd
}

const serverConfigExample = `# ProxyX server configuration — save as <name>.yaml, then:
#   sudo proxyx config apply <name>.yaml
apiVersion: proxyx.io/v1
kind: ProxyConfig
metadata:
  name: my-app          # unique name; becomes <name>.yaml in /etc/proxyx/conf.d
  namespace: default
spec:
  domain: example.com

  # Optional TLS. Omit to serve plain HTTP on :80.
  tls:
    certFile: /etc/letsencrypt/live/example.com/fullchain.pem
    keyFile: /etc/letsencrypt/live/example.com/privkey.pem

  # Optional per-domain rate limit (per client IP). Omit to use the global
  # default from /etc/proxyx/config/config.yaml.
  rateLimit:
    requests: 100
    windowSeconds: 1

  routes:
    # Reverse proxy (round-robin load balancing across servers)
    - name: api
      path: /api/**
      type: ReverseProxy
      reverseProxy:
        servers:
          - url: http://127.0.0.1:8080
          - url: http://127.0.0.1:8081

    # Static file serving
    - name: assets
      path: /static/**
      type: Static
      static:
        root: /var/www/example.com

    # WebSocket passthrough
    - name: realtime
      path: /ws/**
      type: Websocket
      websocket:
        url: ws://127.0.0.1:9000
`

const globalConfigExample = `# ProxyX global configuration — /etc/proxyx/config/config.yaml
# Applies to the whole proxy. Created with defaults on first run; edit and
# restart ProxyX to take effect ( sudo proxyx restart ).
http:
  readTimeout: 3s
  readHeaderTimeout: 2s
  writeTimeout: 5s
  idleTimeout: 15s
  maxHeaderBytes: 1048576
https:
  readTimeout: 5s
  readHeaderTimeout: 3s
  writeTimeout: 10s
  idleTimeout: 30s
  maxHeaderBytes: 1048576
healthCheck:
  enabled: true
  path: /health
  interval: 10s
certbot:
  email: you@example.com

# Default rate limit for any domain whose own spec sets no rateLimit.
rateLimit:
  requests: 1200
  windowSeconds: 1

# Global IP whitelist.
#   mode: bypass  -> listed IPs skip rate limiting; everyone else is served
#   mode: allow   -> only listed IPs may connect; all others get 403
# Entries may be single IPs or CIDR blocks. An empty list disables the whitelist.
whitelist:
  mode: bypass
  ips:
    - 127.0.0.1
    - 10.0.0.0/8
`
