package proxy

import (
	"ProxyX/internal/common"
	"net/http"
	"sort"
	"time"
)

type routeInfo struct {
		lb   *LoadBalancer
		rl   *RateLimiter
		rf   *common.RouteConfig
	}

func NewRouter(config *common.ProxyConfig) http.Handler {
	servers := make(map[string][]routeInfo)

	for _, server := range config.Servers {
		if server.Domain == "" {
			panic("Domain must be specified ")
		}

		rl := NewRateLimiter(server.RateLimit, time.Duration(server.RateWindow)*time.Second)
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

			routes = append(routes, routeInfo{lb: lb, rl: rl, rf: &route})
		}

		sort.Slice(routes, func(i, j int) bool {
			return len(routes[i].rf.Path) > len(routes[j].rf.Path)
		})
		
		servers[server.Domain] = routes
	}

	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		handleRequest(w,r, servers)
	})
}
