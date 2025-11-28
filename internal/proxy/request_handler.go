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

	log.Printf("[REQUEST] %s %s %s from %s", r.Host, r.Method, r.URL.Path, r.RemoteAddr)
	host, _, _ := net.SplitHostPort(r.Host)
	routes, ok := servers[host]
	if !ok {
		ServerDefaultPage(w)
		return
	}

	var matched *routeInfo
	for _, rt := range routes {
		routePrefix := rt.rf.Path
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

			// optional: catch-all route
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
	if !matched.rl.Allow(ip) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("429 Too Many Requests"))
		return
	}

	switch matched.rf.Type {
	case "proxy":
		target := matched.lb.Next()
		if target == nil {
			ServerDefaultPage(w)
			return
		}
		log.Printf("[%s] %s %s -> %s", ip, r.Method, r.URL.Path, target.String())
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)

	case "static":
		if matched.rf.Dir == "" {
			ServerDefaultPage(w)
			return
		}

		buildPath := "./web"
		staticDir := filepath.Join(buildPath, matched.rf.Dir)
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

		log.Println("Serving React SPA:", indexFile)
		http.ServeFile(w, r, indexFile)

	default:
		ServerDefaultPage(w)
	}
}
