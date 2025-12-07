package proxy

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strings"
)



func handleRequest(w http.ResponseWriter, r *http.Request, servers map[string][]routeInfo) {
	routes, ok := servers[r.Host]
	if !ok {
		ServerDefaultPage(w)
		return
	}

	var matched *routeInfo
	for _, rt := range routes {
		routePrefix := rt.routeConfig.Path
		strippedPath := r.URL.Path

		if strings.HasSuffix(routePrefix, "/**") {
			base := strings.TrimSuffix(routePrefix, "/**")
			if strippedPath == base {
				matched = &rt
				break
			}

			if strings.HasPrefix(strippedPath, base+"/") {
				matched = &rt
				break
			}

			continue
		} else {
			normalizedRoute := strings.TrimSuffix(routePrefix, "/")
			normalizedPath := strings.TrimSuffix(strippedPath, "/")
			if normalizedPath == normalizedRoute {
				matched = &rt
				break
			}

			if normalizedRoute == "" {
				matched = &rt
				break
			}
		}
	}

	if matched == nil {
		ServerDefaultPage(w)
		return
	}

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	if !matched.rateLimiter.Allow(ip) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("429 Too Many Requests"))
		return
	}

	switch matched.routeConfig.Type {
	case "proxy":
		target := matched.loadBalancer.Next()
		if target == nil {
			ServerDefaultPage(w)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)

	case "static":
		if matched.routeConfig.Static == nil || matched.routeConfig.Static.Root == "" {
			ServerDefaultPage(w)
			return
		}
		
		staticDir := filepath.Join(matched.routeConfig.Static.Root)
		requestedFile := filepath.Join(staticDir, r.URL.Path)
		if info, err := os.Stat(requestedFile); err == nil && !info.IsDir() {
			http.ServeFile(w, r, requestedFile)
			return
		}

		indexFile := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(indexFile); os.IsNotExist(err) {
			log.Println("index.html not found at", indexFile)
			ServerDefaultPage(w)
			return
		}

		http.ServeFile(w, r, indexFile)
	default:
		
		ServerDefaultPage(w)
	}
}
