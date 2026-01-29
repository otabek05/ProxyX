package cli

import (
	"ProxyX/internal/common"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)


func (c *CLI) applyConfigCmd() *cobra.Command {
	return  &cobra.Command{
		Use:   "apply",
		Short: "⚡ Apply a YAML configuration to ProxyX",
		Args: cobra.ExactArgs(1),
		PreRun: requireRoot,
		Run: func(cmd *cobra.Command, args []string) {
		 filePath := args[0]
		  c.runApply(filePath)
		},
	}
}

func (c *CLI) runApply(file string) {
	if file == "" {
		fmt.Println("❌ Please provide a config file path using -f")
		return
	}

	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("❌ Cannot read file:", err)
		return
	}

	var server common.ServerConfig
	if err := yaml.Unmarshal(data, &server); err != nil {
		fmt.Println("❌ Invalid YAML:", err)
		return
	}

	if err := isValidFormat(&server); err != nil {
		fmt.Println("❌ Invalid configuration format:", err)
		return
	}

	destFile := filepath.Join(
		c.serviceConfig,
		filepath.Base(server.Metadata.Name+".yaml"),
	)

	if err := hasRouteConflict(&server, destFile); err != nil {
		fmt.Println("❌ Route conflict detected:", err)
		return
	}

	if err := os.MkdirAll(c.serviceConfig, 0755); err != nil {
		fmt.Println("❌ Failed to create config directory:", err)
		return
	}

	if server.Spec.RateLimit == nil {
		server.Spec.RateLimit = &common.RateLimitConfig{
			Requests:      1200,
			WindowSeconds: 1,
		}
	}

	if err := os.WriteFile(destFile, data, 0644); err != nil {
		fmt.Println("❌ Failed to write config file:", err)
		return
	}

	fmt.Println("✅ Configuration applied successfully")
	fmt.Println("🔄 Restarting ProxyX service...")

	if err := c.Service.Restart(); err != nil {
		fmt.Println("❌ Failed to restart ProxyX:", err)
		return
	}

	fmt.Println("✅ ProxyX restarted successfully")
}



func isValidFormat(srv *common.ServerConfig) error {
	if srv.Spec.Domain == "" {
		return fmt.Errorf("server missing domain")
	}

	for _, route := range srv.Spec.Routes {
		if route.Path == "" {
			return fmt.Errorf("server '%s' has route missing path", srv.Spec.Domain)
		}

		if !route.Type.IsValid() {
			return fmt.Errorf("server: '%s' has invalid type", srv.Spec.Domain)
		}

		switch route.Type {
		case common.RouteReverseProxy:
			if len(route.ReverseProxy.Servers) == 0 {
				return fmt.Errorf("server '%s' route '%s' of type 'proxy' has no backends", srv.Spec.Domain, route.Path)
			}

		case common.RouteStatic:
			if route.Static == nil || route.Static.Root == "" {
				return fmt.Errorf("server '%s' route '%s' of type 'static' missing dir", srv.Spec.Domain, route.Path)
			}
		case common.RouteWebsocket:
			if route.Websocket == nil || route.Websocket.URL == "" {
				return fmt.Errorf("server '%s' route '%s' of type 'websocket' missing url", srv.Spec.Domain, route.Path)
			}
		default:
			return fmt.Errorf("server '%s' route '%s' has invalid type '%s'", srv.Spec.Domain, route.Path, route.Type)
		}
	}

	return nil
}

func hasRouteConflict(newCfg *common.ServerConfig, newCfgFile string) error {
	configDir := "/etc/proxyx/conf.d"
	files, err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if file == newCfgFile {
			return nil
		}

		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		var existingCfg common.ServerConfig
		if err := yaml.Unmarshal(data, &existingCfg); err != nil {
			return err
		}

		for _, newRoute := range newCfg.Spec.Routes {

			if existingCfg.Spec.Domain != newCfg.Spec.Domain {
				continue
			}

			for _, oldRoute := range existingCfg.Spec.Routes {
				if oldRoute.Path == newRoute.Path {
					return fmt.Errorf(
						"domain='%s' path='%s' already exists in %s",
						newCfg.Spec.Domain,
						newRoute.Path,
						existingCfg.Metadata.Name,
					)
				}
			}
		}
	}

	return nil
}