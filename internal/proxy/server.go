package proxy

import (
	"ProxyX/internal/common"
	"ProxyX/internal/healthchecker"
	"ProxyX/internal/pkg/config"
	"ProxyX/internal/webui"
	"crypto/tls"
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"github.com/valyala/fasthttp"
)

// confDir is where server configurations live (shared with the CLI).
const confDir = "/etc/proxyx/conf.d"

type ProxyServer struct {
	router      atomic.Value // fasthttp.RequestHandler, swapped on reload
	configCount atomic.Int32 // number of loaded server configs
	proxyConfig *common.ProxyConfig
	config      []common.ServerConfig // guarded by reloadMu
	whitelist   *Whitelist

	reloadMu  sync.Mutex
	certCache sync.Map
	proxies   sync.Map
	wsProxies sync.Map
}

func NewServer(config []common.ServerConfig, proxyConfig *common.ProxyConfig) *ProxyServer {
	p := &ProxyServer{
		config:      config,
		proxyConfig: proxyConfig,
		whitelist:   NewWhitelist(proxyConfig.Whitelist),
	}

	p.configureProxy()
	p.configureWSProxy()
	return p
}

// dispatch routes a request through the current (atomically swappable) router.
func (p *ProxyServer) dispatch(ctx *fasthttp.RequestCtx) {
	if h, ok := p.router.Load().(fasthttp.RequestHandler); ok && h != nil {
		h(ctx)
		return
	}
	ServeDefault(ctx)
}

func (p *ProxyServer) Start() error {
	if err := p.loadCertsFor(p.config); err != nil {
		return err
	}

	if p.proxyConfig.HealthCheck.Enabled {
		healthchecker.Start(p.proxyConfig.HealthCheck.Interval)
	}

	// Localhost-only browser config panel, in the same process as the proxy.
	p.startWebPanel()

	go p.runHTTP()

	tlsConfig := &tls.Config{
		GetCertificate: p.getCertificate,
		MinVersion:     tls.VersionTLS12,
	}

	httpsServer := &fasthttp.Server{
		Handler:            p.dispatch,
		TLSConfig:          tlsConfig,
		ReadTimeout:        p.proxyConfig.HTTPS.ReadTimeout,
		WriteTimeout:       p.proxyConfig.HTTPS.WriteTimeout,
		IdleTimeout:        p.proxyConfig.HTTPS.IdleTimeout,
		ReadBufferSize:     32 * 1024,
		WriteBufferSize:    32 * 1024,
		MaxRequestBodySize: 1 * 1024 * 1024,

		DisableHeaderNamesNormalizing: true,
		DisableKeepalive:              false,
		DisablePreParseMultipartForm:  true,
		NoDefaultServerHeader:         true,
		NoDefaultDate:                 true,

		Concurrency: 262144,
	}

	log.Println("HTTPS Proxy server running on :443")
	return httpsServer.ListenAndServeTLS(":443", "", "")
}

// startWebPanel launches the localhost config panel in a goroutine. The panel
// is on by default: it is skipped only when config.yaml explicitly sets
// webPanel.enabled to false. An absent webPanel block (e.g. a config.yaml
// written by an older version) still enables it.
func (p *ProxyServer) startWebPanel() {
	wp := p.proxyConfig.WebPanel
	if wp != nil && !wp.Enabled {
		return // explicitly disabled
	}
	port := 8787
	if wp != nil && wp.Port > 0 && wp.Port <= 65535 {
		port = wp.Port
	}
	panel := &webui.Server{Dir: confDir, Restart: p.Reload}
	go func() {
		log.Printf("Config panel running on http://127.0.0.1:%d (localhost only)", port)
		if err := panel.ListenAndServe(port); err != nil {
			log.Printf("config panel stopped: %v", err)
		}
	}()
}

// Reload re-reads the configurations from disk and swaps the router in place,
// without restarting the process. The browser config panel calls this after a
// change. On any error the previous configuration stays active.
func (p *ProxyServer) Reload() error {
	p.reloadMu.Lock()
	defer p.reloadMu.Unlock()

	cfgs, err := config.LoadConfig()
	if err != nil {
		return err
	}
	handler, err := p.buildRouter(cfgs)
	if err != nil {
		return err
	}
	if err := p.loadCertsFor(cfgs); err != nil {
		return err
	}
	p.config = cfgs
	p.configureWSProxy()
	p.router.Store(handler)
	p.configCount.Store(int32(len(cfgs)))
	log.Printf("configuration reloaded: %d server(s)", len(cfgs))
	return nil
}

// loadCertsFor loads TLS certificates for the given configs into the cache.
func (p *ProxyServer) loadCertsFor(cfgs []common.ServerConfig) error {
	for _, srv := range cfgs {
		if srv.Spec.TLS == nil {
			continue
		}
		cert, err := tls.LoadX509KeyPair(srv.Spec.TLS.CertFile, srv.Spec.TLS.KeyFile)
		if err != nil {
			log.Printf("TLS load failed for %s: %v", srv.Spec.Domain, err)
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
		if p.configCount.Load() == 0 {
			ctx.SetStatusCode(fasthttp.StatusOK)
			ServeDefault(ctx)
			return
		}

		if _, ok := p.certCache.Load(string(ctx.Host())); ok {
			target := "https://" + string(ctx.Host()) + string(ctx.RequestURI())
			ctx.Redirect(target, fasthttp.StatusMovedPermanently)
			return
		}

		p.dispatch(ctx)
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
