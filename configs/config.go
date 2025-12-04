package configs

import (
	"ProxyX/internal/common"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)


func LoadConfig() (*common.ProxyConfig, error) {
	configDir := "/etc/proxyx/configs"
	finalConfig := &common.ProxyConfig{}

	files , err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
	if err != nil {
		return  nil, err 
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return  nil, err
		}

		var cfg common.ProxyConfig
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil,  err
		}

		mergeConfigs(finalConfig, &cfg)
	}


	return finalConfig, nil
}



func mergeConfigs(dst, src *common.ProxyConfig) {
	for _,srcServer := range src.Servers {
		found := false
		for i := range dst.Servers {
			if strings.EqualFold(dst.Servers[i].Domain, srcServer.Domain) {
				dst.Servers[i].Routes = append(dst.Servers[i].Routes, srcServer.Routes...)
				found = true
				break
			}
		}

		if !found {
			dst.Servers = append(dst.Servers, srcServer)
		}
	}
}