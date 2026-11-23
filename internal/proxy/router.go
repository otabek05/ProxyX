package proxy

import (
	"ProxyX/internal/common"
	"sort"
	"time"

	"github.com/valyala/fasthttp"
)

type routeInfo struct {
	loadBalancer *LoadBalancer
	rateLimiter  *RateLimiter
	routeConfig  *common.RouteConfig
}

func (p *ProxyServer) configureProxy()  {
	servers := make(map[string][]routeInfo)

	for _, server := range p.config {
		if server.Spec.Domain == "" {
			panic("Domain must be specified ")
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
					panic(err)
				}
			}

			routes = append(routes, routeInfo{loadBalancer: lb, rateLimiter: rl, routeConfig: &route})
		}

		sort.Slice(routes, func(i, j int) bool {
			return len(routes[i].routeConfig.Path) > len(routes[j].routeConfig.Path)
		})

		servers[server.Spec.Domain] = routes
	}

	p.router = func(ctx *fasthttp.RequestCtx) {
		p.handleRequest(ctx, servers)
	}
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
