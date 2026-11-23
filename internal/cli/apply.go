package cli

import (
	"ProxyX/internal/common"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func (c *CLI) applyConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "apply [file]",
		Short:   "⚡ Apply a YAML configuration to ProxyX",
		Args:    cobra.ExactArgs(1),
		PreRun:  requireRoot,
		Example: "  sudo proxyx config apply ./my-app.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.runApply(args[0]); err != nil {
				color.Red.Println(err.Error())
				os.Exit(1)
			}
		},
	}
}


func (c *CLI) runApply(file string) error {
	if file == "" {
		return fmt.Errorf("❌ please provide a config file path, e.g. proxyx config apply ./my-app.yaml")
	}

	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("❌ cannot read file: %w", err)
	}

	var server common.ServerConfig
	if err := yaml.Unmarshal(data, &server); err != nil {
		return fmt.Errorf("❌ invalid YAML: %w", err)
	}

	if err := isValidFormat(&server); err != nil {
		return fmt.Errorf("❌ invalid configuration format: %w", err)
	}

	destFile := filepath.Join(
		c.serviceConfig,
		server.Metadata.Name+".yaml",
	)

	if err := hasRouteConflict(&server, destFile); err != nil {
		return fmt.Errorf("❌ route conflict detected: %w", err)
	}

	if err := os.MkdirAll(c.serviceConfig, 0755); err != nil {
		return fmt.Errorf("❌ failed to create config directory: %w", err)
	}

	// A config without its own rateLimit falls back to the global default at
	// runtime, so nothing extra needs to be written here.
	if err := os.WriteFile(destFile, data, 0644); err != nil {
		return fmt.Errorf("❌ failed to write config file: %w", err)
	}

	color.Green.Println("✅ Configuration applied successfully")
	fmt.Println("🔄 Restarting ProxyX service...")

	if err := c.Service.Restart(); err != nil {
		return fmt.Errorf("❌ failed to restart ProxyX: %w", err)
	}

	color.Green.Println("✅ ProxyX restarted successfully")

	return nil
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
