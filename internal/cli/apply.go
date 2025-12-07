package cli

import (
	"ProxyX/internal/common"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var applyFile string

func init() {
	rootCmd.AddCommand(applyCmd)
	applyCmd.Flags().StringVarP(&applyFile, "file", "f", "", "Path to config file to add")
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply configuration file to ProxyX",
	Run: func(cmd *cobra.Command, args []string) {
		if applyFile == "" {
			fmt.Println("Please provide a config file path using -f")
			return
		}

		data, err := os.ReadFile(applyFile)
		if err != nil {
			fmt.Println("Cannot read file:", err)
			return
		}

		var server common.ServerConfig
		err = yaml.Unmarshal(data, &server)
		if err != nil {
			fmt.Println("Invalid YAML:", err)
			return
		}

		if err := isValidFormat(&server); err != nil {
			fmt.Println(err.Error())
			return
		}

		destDir := "/etc/proxyx/configs"
		desFile := filepath.Join(destDir, filepath.Base(applyFile))
		if err := hasRouteConflict(&server, desFile); err != nil {
			fmt.Println(err.Error())
			return
		}

		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			fmt.Println("Failed to created dir: ", err)
			return
		}

		if server.Spec.RateLimit == nil {
			server.Spec.RateLimit = &common.RateLimitConfig{
				Requests: 1200,
				WindowSeconds: 1,
			}
		}

		err = os.WriteFile(desFile, data, 0644)
		if err != nil {
			fmt.Println("Failed to write config file:", err)
			return
		}

		fmt.Println("Configuration applied successfully")
		reloadProxyX()
	},
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
			return  fmt.Errorf("server: '%s' has invalid type", srv.Spec.Domain)
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

		default:
			return fmt.Errorf("server '%s' route '%s' has invalid type '%s'", srv.Spec.Domain, route.Path, route.Type)
		}
	}

	return nil
}

func hasRouteConflict(newCfg *common.ServerConfig, newCfgFile string) error {
	configDir := "/etc/proxyx/configs"
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
						"conflict detected: domain='%s' path='%s' already exists in %s",
						newCfg.Spec.Domain,
						newRoute.Path,
						filepath.Base(file),
					)
				}
			}
		}
	}

	return nil
}

func hasDuplicate[T comparable](slice []T) error {
	seen := make(map[T]bool)
	for _, v := range slice {
		if seen[v] {
			return fmt.Errorf("duplicate: %v", v)
		}

		seen[v] = true
	}

	return nil
}
