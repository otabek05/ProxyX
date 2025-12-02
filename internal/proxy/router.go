package proxy

import (
	"ProxyX/internal/common"
	"net/http"
	"sort"
	"time"
)

type routeInfo struct {
		loadBalancer   *LoadBalancer
		rateLimiter   *RateLimiter
		routeConfig   *common.RouteConfig
	}

func NewRouter(config *common.ProxyConfig) http.Handler {
	servers := make(map[string][]routeInfo)

	for _, server := range config.Servers {
		if server.Domain == "" {
			panic("Domain must be specified ")
		}

		//rl := NewRateLimiter(server.RateLimit, time.Duration(server.RateWindow)*time.Second)
		var routes []routeInfo
		for _, route := range server.Routes {

		    if route.Type == "" {
				route.Type = "proxy"
			}

			var lb *LoadBalancer
			if route.Type == "proxy" {
				var err error
				lb, err = NewLoadBalancer(route.Backends)
				if err != nil {
					panic(err)
				}
			}

			rl := NewRateLimiter(route.RateLimit, time.Duration(route.RateWindow)*time.Second)

			routes = append(routes, routeInfo{loadBalancer: lb, rateLimiter: rl, routeConfig: &route})
		}

		sort.Slice(routes, func(i, j int) bool {
			return len(routes[i].routeConfig.Path) > len(routes[j].routeConfig.Path)
		})
		
		servers[server.Domain] = routes
	}

	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		handleRequest(w,r, servers)
	})
}
