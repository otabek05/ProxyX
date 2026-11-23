package common

import "time"

type ProxyConfig struct {
	HTTP        HTTPConfig        `yaml:"http"`
	HTTPS       HTTPSConfig       `yaml:"https"`
	HealthCheck HealthCheckConfig `yaml:"healthCheck"`
	Certbot     Certbot           `yaml:"certbot"`

	// RateLimit is the default per-IP rate limit applied to any domain whose
	// own spec does not define one.
	RateLimit *RateLimitConfig `yaml:"rateLimit"`

	// Whitelist controls global IP access: either exempting trusted IPs from
	// rate limiting (mode: bypass) or restricting access to them (mode: allow).
	Whitelist *WhitelistConfig `yaml:"whitelist"`
}

// Whitelist modes.
const (
	WhitelistBypass = "bypass" // listed IPs skip rate limiting; everyone else is served normally
	WhitelistAllow  = "allow"  // only listed IPs may connect; all others get 403
)

type WhitelistConfig struct {
	Mode string   `yaml:"mode"` // "bypass" (default) or "allow"
	IPs  []string `yaml:"ips"`  // individual IPs or CIDR blocks, e.g. 10.0.0.5, 203.0.113.0/24
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
