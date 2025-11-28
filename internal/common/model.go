package common


type RouteConfig struct {
	Path  string `yaml:"path"`
	Type  string `yaml:"type,omitempty"`
	Dir   string  `yaml:"dir,omitempty"`
	LoadBalancer string `yaml:"loadbalancer"`
	Backends   []string `yaml:"backends"`
}

type ServerConfig  struct {
	Domain  string `yaml:"domain"`
	RateLimit int  `yaml:"rate_limit"`
	RateWindow int  `yaml:"rate_window"`
	Routes  []RouteConfig `yaml:"routes"`
}

type ProxyConfig struct {
	Servers  []ServerConfig `yaml:"servers"`
}


