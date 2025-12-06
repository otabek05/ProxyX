package configs

import (
	"ProxyX/internal/common"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)


func LoadConfig() ([]common.ServerConfig, error)  {
	configDir := "/etc/proxyx/configs"
	var finalConfig []common.ServerConfig

	files , err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
	if err != nil {
		return  nil, err 
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return  nil, err
		}

		var cfg common.ServerConfig
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil,  err
		}

		finalConfig = mergeConfigs(finalConfig, &cfg)
	}


	return finalConfig, nil
}



func mergeConfigs(dst []common.ServerConfig, src *common.ServerConfig) []common.ServerConfig {
	    found := false
		for i := range dst {
			if strings.EqualFold(dst[i].Domain, src.Domain) {
				dst[i].Routes = append(dst[i].Routes, src.Routes...)
				found = true
				break
			}
		}

		if !found {
			dst = append(dst, *src)
		}

		return dst

}