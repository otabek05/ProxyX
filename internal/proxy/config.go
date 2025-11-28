package proxy

import (
	"ProxyX/internal/common"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func LoadConfig(dir string) (*common.ProxyConfig, error) {
	finalConfig := &common.ProxyConfig{}

	files , err := filepath.Glob(filepath.Join(dir, "*.yaml"))
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
	dst.Servers = append(dst.Servers, src.Servers...)
}
