package common

import "time"

type ProxyConfig struct {
	HTTP        HTTPConfig        `yaml:"http"`
	HTTPS       HTTPSConfig       `yaml:"https"`
	HealthCheck HealthCheckConfig `yaml:"healthCheck"`
	Certbot     Certbot           `yaml:"certbot"`
}

type Certbot struct {
	Email string `yaml:"email"`
}

type HTTPConfig struct {
	ReadTimeout       time.Duration `yaml:"readTimeout"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
	WriteTimeout      time.Duration `yaml:"writeTimeout"`
	IdleTimeout       time.Duration `yaml:"idleTimeout"`
	MaxHeaderBytes    int           `yaml:"maxHeaderBytes"`
}

type HTTPSConfig struct {
	ReadTimeout       time.Duration `yaml:"readTimeout"`
	ReadHeaderTimeout time.Duration `yaml:"readHeaderTimeout"`
	WriteTimeout      time.Duration `yaml:"writeTimeout"`
	IdleTimeout       time.Duration `yaml:"idleTimeout"`
	MaxHeaderBytes    int           `yaml:"maxHeaderBytes"`
}

type HealthCheckConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Path     string        `yaml:"path"`
	Interval time.Duration `yaml:"interval"`
}
