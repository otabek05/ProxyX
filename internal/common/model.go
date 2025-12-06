package common


type RouteConfig struct {
	Path  string `yaml:"path"`
	Type  string `yaml:"type,omitempty"`
	Dir   string  `yaml:"dir,omitempty"`
	Backends   []string `yaml:"backends"`
}

type ServerConfig  struct {
	Domain  string `yaml:"domain"`
	CertFile string `yaml:"cert_file,omitempty"`
	KeyFile string `yaml:"key_file,omitempty"`
	RateLimit int `yaml:"rate_limit,omitempty"`
	RateWindow int `yaml:"rate_window,omitempty"`
	Routes  []RouteConfig `yaml:"routes"`
}


/*
type ProxyConfig struct {
	Servers  []ServerConfig `yaml:"servers"`
}

*/


func  ToDomainList(p []ServerConfig) []string {
	var domains []string
	for _, server := range p{
		domains = append(domains, server.Domain)
	}

	return domains
}


/*
func (r *RouteConfig) ApplyDefaults() {
    if r.RateLimit == 0 {
        r.RateLimit = 100   
    }
    if r.RateWindow == 0 {
        r.RateWindow = 60  
    }
}
	*/
