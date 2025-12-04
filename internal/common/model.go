package common


type RouteConfig struct {
	Path  string `yaml:"path"`
	Type  string `yaml:"type,omitempty"`
	RateLimit int `yaml:"rate_limit,omitempty"`
	RateWindow int `yaml:"rate_window,omitempty"`
	Dir   string  `yaml:"dir,omitempty"`
	Backends   []string `yaml:"backends"`
}

type ServerConfig  struct {
	Domain  string `yaml:"domain"`
	Routes  []RouteConfig `yaml:"routes"`
}



type ProxyConfig struct {
	Servers  []ServerConfig `yaml:"servers"`
}


func (p *ProxyConfig) ToDomainList() []string {
	var domains []string
	for _, server := range p.Servers{
		domains = append(domains, server.Domain)
	}

	return domains
}



func (r *RouteConfig) ApplyDefaults() {
    if r.RateLimit == 0 {
        r.RateLimit = 100   
    }
    if r.RateWindow == 0 {
        r.RateWindow = 60  
    }
}
