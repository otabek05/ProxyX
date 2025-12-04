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

		var cfg common.ProxyConfig
		err = yaml.Unmarshal(data, &cfg)
		if err != nil {
			fmt.Println("Invalid YAML:", err)
			return
		}

		if err := isValidFormat(&cfg); err != nil {
			fmt.Println(err.Error())
			return
		}

		destDir := "/etc/proxyx/configs"
		desFile := filepath.Join(destDir, filepath.Base(applyFile))
		if err := hasRouteConflict(&cfg, desFile); err != nil {
			fmt.Println(err.Error())
			return
		}

		
		err = os.MkdirAll(destDir, 0755)
		if err != nil {
			fmt.Println("Failed to created dir: ", err)
			return
		}

		err = os.WriteFile(desFile, data, 0644)
		if err != nil {
			fmt.Println("Failed to write config file:", err)
			return
		}

		fmt.Printf("Configuration applied successfully: %s\n", desFile)
		reloadProxyX()
	},
}

func isValidFormat(config *common.ProxyConfig) error {
	if len(config.Servers) == 0 {
		return fmt.Errorf("no servers defined in config")
	}

	for _, srv := range config.Servers {
		if srv.Domain == "" {
			return fmt.Errorf("server missing domain")
		}
		fmt.Printf("Server: %s\n", srv.Domain)

		for _, route := range srv.Routes {
			if route.Path == "" {
				return fmt.Errorf("server '%s' has route missing path", srv.Domain)
			}
			fmt.Printf("  Route path: %s\n", route.Path)

			switch route.Type {
			case "proxy":
				if len(route.Backends) == 0 {
					return fmt.Errorf("server '%s' route '%s' of type 'proxy' has no backends", srv.Domain, route.Path)
				}
				fmt.Printf("    Type: proxy, Backends: %v\n", route.Backends)
			case "static":
				if route.Dir == "" {
					return fmt.Errorf("server '%s' route '%s' of type 'static' missing dir", srv.Domain, route.Path)
				}
				fmt.Printf("    Type: static, Dir: %s\n", route.Dir)
			default:
				return fmt.Errorf("server '%s' route '%s' has invalid type '%s'", srv.Domain, route.Path, route.Type)
			}
		}
	}

	return nil
}

func hasRouteConflict(newCfg *common.ProxyConfig, newCfgFile string ) error {
	configDir := "/etc/proxyx/configs"
	files, err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println("File: ", file)
		if file == newCfgFile {
			return nil
		}

		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		var existingCfg common.ProxyConfig
		if err := yaml.Unmarshal(data, &existingCfg); err != nil {
			return err
		}

		for _, newServer := range newCfg.Servers {
			for _, newRoute := range newServer.Routes {

				for _, oldServer := range existingCfg.Servers {
					if oldServer.Domain != newServer.Domain {
						continue
					}

					for _, oldRoute := range oldServer.Routes {
						if oldRoute.Path == newRoute.Path {
							return fmt.Errorf(
								"conflict detected: domain='%s' path='%s' already exists in %s",
								newServer.Domain,
								newRoute.Path,
								filepath.Base(file),
							)
						}
					}
				}
			}
		}
	}

	return nil
}
