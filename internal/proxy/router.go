package proxy

import (
	"ProxyX/internal/common"
	"fmt"
	"sort"
	"time"

	"github.com/valyala/fasthttp"
)

type routeInfo struct {
	loadBalancer *LoadBalancer
	rateLimiter  *RateLimiter
	routeConfig  *common.RouteConfig
}

// configureProxy builds the initial router and stores it. It panics on a bad
// configuration at startup (preserving the original fail-fast behavior).
func (p *ProxyServer) configureProxy() {
	handler, err := p.buildRouter(p.config)
	if err != nil {
		panic(err)
	}
	p.router.Store(handler)
	p.configCount.Store(int32(len(p.config)))
}

// buildRouter builds a request handler for the given configs. Unlike the
// startup path it returns an error instead of panicking, so a reload triggered
// at runtime can fail safely and keep the previous router.
func (p *ProxyServer) buildRouter(configs []common.ServerConfig) (fasthttp.RequestHandler, error) {
	servers := make(map[string][]routeInfo)

	for _, server := range configs {
		if server.Spec.Domain == "" {
			return nil, fmt.Errorf("a configuration is missing its domain")
		}

		rlCfg := p.rateLimitFor(&server)
		rl := NewRateLimiter(rlCfg.Requests, time.Duration(rlCfg.WindowSeconds)*time.Second)
		var routes []routeInfo
		for _, route := range server.Spec.Routes {
			if route.Type == "" {
				route.Type = common.RouteReverseProxy
			}

			var lb *LoadBalancer
			if route.Type == common.RouteReverseProxy {
				var err error
				lb, err = NewLoadBalancer(route.ReverseProxy.Servers, p.proxyConfig)
				if err != nil {
					return nil, err
				}
			}

			routes = append(routes, routeInfo{loadBalancer: lb, rateLimiter: rl, routeConfig: &route})
		}

		sort.Slice(routes, func(i, j int) bool {
			return len(routes[i].routeConfig.Path) > len(routes[j].routeConfig.Path)
		})

		servers[server.Spec.Domain] = routes
	}

	return func(ctx *fasthttp.RequestCtx) {
		p.handleRequest(ctx, servers)
	}, nil
}

// rateLimitFor resolves the rate limit for a server: its own spec wins, else the
// global default from config.yaml, else a hardcoded safe fallback. This also
// guards against a nil per-domain rate limit (a config applied without one).
func (p *ProxyServer) rateLimitFor(server *common.ServerConfig) common.RateLimitConfig {
	if server.Spec.RateLimit != nil {
		return *server.Spec.RateLimit
	}
	if p.proxyConfig != nil && p.proxyConfig.RateLimit != nil {
		return *p.proxyConfig.RateLimit
	}
	return common.RateLimitConfig{Requests: 1200, WindowSeconds: 1}
}
