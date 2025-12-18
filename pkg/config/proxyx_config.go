package config

import (
	"ProxyX/internal/common"
	"errors"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)


func LoadProxyXConfig() (*common.ProxyConfig, error ) {
	path := "/etc/proxyx/config/config.yaml"

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if err := createDefaultProxyConfig(path); err != nil {
			return nil, err 
		}
	} 

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err 
	}

	var config common.ProxyConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
	  return nil, err 
	}

	return &config, nil
}


func createDefaultProxyConfig(path string ) error  {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err 
	}

	cfg := defaultConfig()
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err 
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err 
	}

	return os.Rename(tmp, path)
}


func defaultConfig() common.ProxyConfig {
	return common.ProxyConfig{
		HTTP: common.HTTPConfig{
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 3 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       20 * time.Second,
			MaxHeaderBytes:    1024 * 1024, // 1MB
		},
		HTTPS: common.HTTPSConfig{
			ReadTimeout:       15 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
			MaxHeaderBytes:    1024 * 1024, // 1MB
		},
		HealthCheck: common.HealthCheckConfig{
			Enabled:  true,
			Path:     "/health",
			Interval: 10 * time.Second,
		},
	}
}