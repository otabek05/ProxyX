package proxy

import (
	"ProxyX/internal/common"
	"strings"

	"github.com/valyala/fasthttp"
)

func (p *ProxyServer) handleRequest(ctx *fasthttp.RequestCtx, servers map[string][]routeInfo) {
	host := string(ctx.Host())
	path := string(ctx.Path())

	routes, ok := servers[host]
	if !ok {
		ServeDefault(ctx)
		return
	}

	matched := findMatchingRoute(routes, path)
	if matched == nil {
		ServeError(ctx)
		return
	}

	ip := ctx.RemoteAddr().String()
    if !matched.rateLimiter.Allow(ip) {
        ctx.SetStatusCode(fasthttp.StatusTooManyRequests)
        ctx.SetBodyString("429 Too Many Requests")
        return
    }
	

	switch matched.routeConfig.Type {
	case common.RouteReverseProxy:
		p.reverseProxyxHandler(ctx, matched)
	case common.RouteStatic:
		staticRouteHandler(ctx, matched)
	case common.RouteWebsocket:
		p.websocketProxyHandler(ctx)
	default:
		ServeError(ctx)
	}
}

func findMatchingRoute(routes []routeInfo, path string) *routeInfo {
	for _, rt := range routes {
		routePrefix := rt.routeConfig.Path
		strippedPath := path

		if strings.HasSuffix(routePrefix, "/**") {
			base := strings.TrimSuffix(routePrefix, "/**")
			if strippedPath == base {
				return &rt
			}

			if strings.HasPrefix(strippedPath, base+"/") {
				return &rt
			}

			continue
		} else {
			normalizedRoute := strings.TrimSuffix(routePrefix, "/")
			normalizedPath := strings.TrimSuffix(strippedPath, "/")
			if normalizedPath == normalizedRoute {
				return &rt
			}

			if normalizedRoute == "" {
				return &rt
			}
		}
	}

	return nil
}

