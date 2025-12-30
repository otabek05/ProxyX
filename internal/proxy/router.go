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

func (p *ProxyServer) NewRouter(config []common.ServerConfig, proxyConfig *common.ProxyConfig) fasthttp.RequestHandler {
	servers := make(map[string][]routeInfo)

	for _, server := range config {
		if server.Spec.Domain == "" {
			panic("Domain must be specified ")
		}

		rl := NewRateLimiter(server.Spec.RateLimit.Requests, time.Duration(server.Spec.RateLimit.WindowSeconds)*time.Minute)
		var routes []routeInfo
		for _, route := range server.Spec.Routes {

			if route.Type == "" {
				route.Type = common.RouteReverseProxy
			}

			var lb *LoadBalancer
			if route.Type == common.RouteReverseProxy {
				var err error
				lb, err = NewLoadBalancer(route.ReverseProxy.Servers, proxyConfig)
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

	return func(ctx *fasthttp.RequestCtx) {
		p.handleRequest(ctx, servers)
	}
}
