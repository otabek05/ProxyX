package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func (p *ProxyServer) reverseProxyxHandler(w http.ResponseWriter, r *http.Request, matched *routeInfo) {
	target := matched.loadBalancer.Next()
	if target == nil {
		http.Error(w, "No upstream", http.StatusServiceUnavailable)
		return
	}

	key := target.Scheme + "://" + target.Host
	proxy, ok := p.proxies[key]
	if !ok {
		u := &url.URL{
			Scheme: target.Scheme,
			Host:   target.Host,
		}
		proxy = newUpstreamProxy(u)
		p.proxies[key] = proxy
	}

	r.Header.Set("X-Forwarded-For", r.RemoteAddr)
	r.Header.Set("X-Forwarded-Host", r.Host)
	if r.TLS != nil {
		r.Header.Set("X-Forwarded-Proto", "https")
	} else {
		r.Header.Set("X-Forwarded-Proto", "http")
	}

	proxy.ServeHTTP(w, r)
}

func newUpstreamProxy(target *url.URL) *httputil.ReverseProxy {
	transport := &http.Transport{
		MaxIdleConns:        2048,
		MaxIdleConnsPerHost: 2048,
		IdleConnTimeout:    30 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		ResponseHeaderTimeout: 15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression: true,
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = transport

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Println("upstream error:", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	return proxy
}
