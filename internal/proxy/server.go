package proxy

import (
	"ProxyX/internal/common"
	"ProxyX/internal/healthchecker"
	"crypto/tls"
	"fmt"
	"log"
	"github.com/valyala/fasthttp"
)

type ProxyServer struct {
	router      fasthttp.RequestHandler
	proxyConfig *common.ProxyConfig
	config      []common.ServerConfig
	certCache   map[string]*tls.Certificate
	proxies     map[string]*fasthttp.Client
}

func NewServer(config []common.ServerConfig, proxyConfig *common.ProxyConfig) *ProxyServer {
	p := &ProxyServer{
		config: config, 
		proxyConfig: proxyConfig,
		proxies: make(map[string]*fasthttp.Client),
	}
	p.router = p.NewRouter(config, p.proxyConfig)
	return p
}

func (p *ProxyServer) Start() {
	if err := p.loadAllCertificates(); err != nil {
		log.Fatal(err)
	}

	if p.proxyConfig.HealthCheck.Enabled {
		healthchecker.Start(p.proxyConfig.HealthCheck.Interval)
	}



	go p.runHTTP()

	tlsConfig := &tls.Config{
		GetCertificate: p.getCertificate,
	}


	log.Println("HTTPS Proxy server running on :443")
	httpsServer := &fasthttp.Server{
		Handler:            p.router,
		TLSConfig:          tlsConfig,
		ReadTimeout:        p.proxyConfig.HTTPS.ReadTimeout,
		WriteTimeout:       p.proxyConfig.HTTPS.WriteTimeout,
		IdleTimeout:        p.proxyConfig.HTTPS.IdleTimeout,
		ReadBufferSize:     4 *1024 *1024,
		WriteBufferSize:    4 *1024 *1024,
		MaxRequestBodySize: 1024 * 1024,
	}

	log.Println("HTTPS Proxy server running on :443")
	log.Fatal(httpsServer.ListenAndServeTLS(":443", "",""))
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



func (p *ProxyServer) runHTTP() {
	handler := func(ctx *fasthttp.RequestCtx) {
		if len(p.config) == 0 {
			ctx.SetStatusCode(fasthttp.StatusOK)
		    ServeProxyHomepage(ctx)
			return
		}

		if _, ok := p.certCache[string(ctx.Host())]; ok {
			target := "https://" + string(ctx.Host()) + string(ctx.RequestURI())
			ctx.Redirect(target, fasthttp.StatusMovedPermanently)
			return
		}

		p.router(ctx)
	}

	log.Println("HTTP Proxy server running on :80")
	server := &fasthttp.Server{
		Handler:         handler,
		ReadTimeout:     p.proxyConfig.HTTP.ReadTimeout,
		WriteTimeout:    p.proxyConfig.HTTP.WriteTimeout,
		IdleTimeout:     p.proxyConfig.HTTP.IdleTimeout,
		ReadBufferSize:  4 * 1024 *1024,
		WriteBufferSize: 4 * 1024 *1024,
		MaxRequestBodySize: 1024 * 1024,
	}

	log.Fatal(server.ListenAndServe(":80"))
}