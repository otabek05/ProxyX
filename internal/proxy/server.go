package proxy

import (
	"ProxyX/internal/common"
	"ProxyX/internal/healthchecker"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type ProxyServer struct {
	router      http.Handler
	proxyConfig *common.ProxyConfig
	config      []common.ServerConfig
	certCache   map[string]*tls.Certificate
	proxies     map[string]*httputil.ReverseProxy
}

func NewServer(config []common.ServerConfig, proxyConfig *common.ProxyConfig) *ProxyServer {
	p := &ProxyServer{
		config:      config,
		proxyConfig: proxyConfig,
		certCache:   make(map[string]*tls.Certificate),
		proxies:     make(map[string]*httputil.ReverseProxy),
	}

	p.router = http.HandlerFunc(p.re)
	return p
}

func (p *ProxyServer) Start() {
	if err := p.loadAllCertificates(); err != nil {
		log.Fatal(err)
	}

	p.initProxies()

	if p.proxyConfig.HealthCheck.Enabled {
		healthchecker.Start(p.proxyConfig.HealthCheck.Interval)
	}

	go p.runHTTP()
	p.runHTTPS()
}

func (p *ProxyServer) runHTTPS() {
	tlsConfig := &tls.Config{
		GetCertificate: p.getCertificate,
		MinVersion:     tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:         ":443",
		Handler:      p.router,
		TLSConfig:   tlsConfig,
		ReadTimeout: p.proxyConfig.HTTPS.ReadTimeout,
		WriteTimeout: p.proxyConfig.HTTPS.WriteTimeout,
		IdleTimeout:  p.proxyConfig.HTTPS.IdleTimeout,
	}

	log.Println("HTTPS Proxy server running on :443")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func (p *ProxyServer) runHTTP() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if len(p.config) == 0 {
			w.WriteHeader(http.StatusOK)
			ServeProxyHomepageHTTP(w)
			return
		}

		if _, ok := p.certCache[r.Host]; ok {
			target := "https://" + r.Host + r.RequestURI
			http.Redirect(w, r, target, http.StatusMovedPermanently)
			return
		}

		p.router.ServeHTTP(w, r)
	})

	server := &http.Server{
		Addr:         ":80",
		Handler:      handler,
		ReadTimeout:  p.proxyConfig.HTTP.ReadTimeout,
		WriteTimeout: p.proxyConfig.HTTP.WriteTimeout,
		IdleTimeout:  p.proxyConfig.HTTP.IdleTimeout,
	}

	log.Println("HTTP Proxy server running on :80")
	log.Fatal(server.ListenAndServe())
}
/*
func (p *ProxyServer) routeRequest(w http.ResponseWriter, r *http.Request) {
	proxy, ok := p.proxies[r.Host]
	if !ok {
		http.NotFound(w, r)
		return
	}
	proxy.ServeHTTP(w, r)
}

func (p *ProxyServer) initProxies() {
	transport := &http.Transport{
		MaxIdleConns:        1024,
		MaxIdleConnsPerHost: 128,
		IdleConnTimeout:    90 * time.Second,
		DisableCompression: true,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: p.proxyConfig.HTTP.ReadTimeout,
		ExpectContinueTimeout: 1 * time.Second,
	}

	for _, srv := range p.config {
		u, err := url.Parse(srv.Spec.Backend.URL)
		if err != nil {
			log.Printf("Invalid backend URL %s: %v", srv.Spec.Backend.URL, err)
			continue
		}

		proxy := httputil.NewSingleHostReverseProxy(u)
		proxy.Transport = transport
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Proxy error [%s]: %v", r.Host, err)
			http.Error(w, "Bad Gateway", http.StatusBadGateway)
		}

		p.proxies[srv.Spec.Domain] = proxy
	}
}
	*/

func (p *ProxyServer) loadAllCertificates() error {
	for _, srv := range p.config {
		if srv.Spec.TLS == nil {
			continue
		}

		cert, err := tls.LoadX509KeyPair(
			srv.Spec.TLS.CertFile,
			srv.Spec.TLS.KeyFile,
		)
		if err != nil {
			log.Printf("TLS load failed for %s: %v", srv.Spec.Domain, err)
			continue
		}

		p.certCache[srv.Spec.Domain] = &cert
	}
	return nil
}

func (p *ProxyServer) getCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if cert, ok := p.certCache[hello.ServerName]; ok {
		return cert, nil
	}
	return nil, fmt.Errorf("no TLS cert for domain: %s", hello.ServerName)
}
