package common


type ServerConfig struct {
	ApiVersion string    `yaml:"apiVersion"`
	Kind       string    `yaml:"kind"`
	Metadata   Metadata  `yaml:"metadata"`
	Spec       ProxySpec `yaml:"spec"`
}

type Metadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type ProxySpec struct {
	Domain    string          `yaml:"domain"`
	TLS       *TLSConfig      `yaml:"tls"`
	RateLimit *RateLimitConfig `yaml:"rateLimit"`
	Routes    []RouteConfig `yaml:"routes"`
}

type TLSConfig struct {
	CertFile string `yaml:"certFile"`
	KeyFile  string `yaml:"keyFile"`
}

type RateLimitConfig struct {
	Requests      int `yaml:"requests"`
	WindowSeconds int `yaml:"windowSeconds"`
}

type RouteConfig struct {
	Name         string            `yaml:"name"`
	Path         string            `yaml:"path"`
	Type         RouteType            `yaml:"type,omitempty"`
	Static       *StaticConfig     `yaml:"static,omitempty"`
	ReverseProxy *ReverseProxySpec `yaml:"reverseProxy"`
}

type StaticConfig struct {
	Root string `yaml:"root"`
}

type ReverseProxySpec struct {
	Servers []ProxyServer `yaml:"servers"`
}

type ProxyServer struct {
	URL string `yaml:"url"`
}

func ToDomainList(p []ServerConfig) []string {
	var domains []string
	for _, server := range p {
		domains = append(domains, server.Spec.Domain)
	}

	return domains
}
