package proxy

import (
	"ProxyX/internal/common"
	"ProxyX/internal/healthchecker"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ProxyServer struct {
	router    http.Handler
	config    []common.ServerConfig
	certCache map[string]*tls.Certificate
}

func NewServer(config []common.ServerConfig) *ProxyServer {
	p := &ProxyServer{config: config}
	p.router = NewRouter(config)
	return p
}

func (p *ProxyServer) Start() {
	if err := p.loadAllCertificates(); err != nil {
		log.Fatal(err)
	}

	healthchecker.Start(3 * time.Second)

	httpServer := &http.Server{
		Addr: ":80",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if len(p.config) == 0 {
				ServeProxyHomepage(w)
				return
			}

			if _, ok := p.certCache[r.Host]; ok {
				target := "https://" + r.Host + r.URL.String()
				http.Redirect(w, r, target, http.StatusMovedPermanently)
				return
			}

			p.router.ServeHTTP(w, r)
		}),

		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	go func() {
		log.Println("HTTP Proxy server running on :80")
		log.Fatal(httpServer.ListenAndServe())
	}()

	tlsConfig := &tls.Config{
		GetCertificate: p.getCertificate,
	}

	server := &http.Server{
		Addr:              ":443",
		Handler:           p.router,
		TLSConfig:         tlsConfig,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, 
	}

	log.Println("HTTPS Proxy server running on :443")
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func (p *ProxyServer) loadAllCertificates() error {
	p.certCache = make(map[string]*tls.Certificate)
	for _, srv := range p.config {
		if srv.Spec.TLS == nil {
			continue
		}

		cert, err := tls.LoadX509KeyPair(srv.Spec.TLS.CertFile, srv.Spec.TLS.KeyFile)
		if err != nil {
			fmt.Printf("TLS load failed for %s: %v", srv.Spec.Domain, err)
			continue
		}

		p.certCache[srv.Spec.Domain] = &cert
		log.Println("Loaded TLS for:", srv.Spec.Domain)
	}

	return nil
}

func (p *ProxyServer) getCertificate(tslHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	domain := tslHello.ServerName

	if cert, ok := p.certCache[domain]; ok {
		return cert, nil
	}

	return nil, fmt.Errorf("no TLS cert for domain: %s", domain)
}
