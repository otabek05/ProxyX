package proxy

import (
	"ProxyX/internal/common"
	"ProxyX/internal/healthchecker"
	"crypto/tls"
	"fmt"
	"log"
	"sync"

	"github.com/valyala/fasthttp"
)

type ProxyServer struct {
	router      fasthttp.RequestHandler
	proxyConfig *common.ProxyConfig
	config      []common.ServerConfig

	certCache  sync.Map
	proxies   sync.Map
	wsProxies sync.Map

	stats struct {
		TotalRequests   uint64
		UpstreamErrors  uint64
		WebsocketErrors uint64
	}
}

func NewServer(config []common.ServerConfig, proxyConfig *common.ProxyConfig) *ProxyServer {
	p := &ProxyServer{
		config:      config,
		proxyConfig: proxyConfig,
	}

	p.router = p.NewRouter(config, p.proxyConfig)
	p.configureWSProxy()
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
		MinVersion:     tls.VersionTLS12,
	}

	httpsServer := &fasthttp.Server{
		Handler:            p.router,
		TLSConfig:          tlsConfig,
		ReadTimeout:        p.proxyConfig.HTTPS.ReadTimeout,
		WriteTimeout:       p.proxyConfig.HTTPS.WriteTimeout,
		IdleTimeout:        p.proxyConfig.HTTPS.IdleTimeout,
		ReadBufferSize:     32 * 1024,
		WriteBufferSize:    32 * 1024,
		MaxRequestBodySize: 1 * 1024 * 1024,

		DisableHeaderNamesNormalizing: true,
		DisableKeepalive: false,
		DisablePreParseMultipartForm: true,
		NoDefaultServerHeader:         true,
		NoDefaultDate:                 true,

		Concurrency: 262144,
	}

	log.Println("HTTPS Proxy server running on :443")
	log.Fatal(httpsServer.ListenAndServeTLS(":443", "", ""))
}

func (p *ProxyServer) loadAllCertificates() error {
	for _, srv := range p.config {
		if srv.Spec.TLS == nil {
			continue
		}

		cert, err := tls.LoadX509KeyPair(srv.Spec.TLS.CertFile, srv.Spec.TLS.KeyFile)
		if err != nil {
			fmt.Printf("TLS load failed for %s: %v", srv.Spec.Domain, err)
			continue
		}

		p.certCache.Store(srv.Spec.Domain, &cert)
		log.Println("Loaded TLS for:", srv.Spec.Domain)
	}

	return nil
}

func (p *ProxyServer) getCertificate(tslHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	domain := tslHello.ServerName

	if certAny, ok := p.certCache.Load(domain); ok {
		return certAny.(*tls.Certificate), nil
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

		if _, ok := p.certCache.Load(string(ctx.Host())); ok {
			target := "https://" + string(ctx.Host()) + string(ctx.RequestURI())
			ctx.Redirect(target, fasthttp.StatusMovedPermanently)
			return
		}

		p.router(ctx)
	}

	log.Println("HTTP Proxy server running on :80")
	server := &fasthttp.Server{
		Handler:            handler,
		ReadTimeout:        p.proxyConfig.HTTP.ReadTimeout,
		WriteTimeout:       p.proxyConfig.HTTP.WriteTimeout,
		IdleTimeout:        p.proxyConfig.HTTP.IdleTimeout,
		ReadBufferSize:     32 * 1024,
		WriteBufferSize:    32 * 1024,
		MaxRequestBodySize: 1024 * 1024,
		Concurrency:        0,
	}

	log.Fatal(server.ListenAndServe(":80"))
}
