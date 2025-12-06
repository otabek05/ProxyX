package proxy

import (
	"ProxyX/internal/common"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
)


type ProxyServer struct {
	router http.Handler
	config []common.ServerConfig
	certCache  map[string]*tls.Certificate
}

func NewServer(config []common.ServerConfig) *ProxyServer {
	p := &ProxyServer{config: config,}
	p.router = NewRouter(config)
	return p
}

func (p *ProxyServer) Start()  {
	if err := p.loadAllCertificates(); err != nil {
		log.Fatal(err)
	}

	go func ()  {
	   log.Println("HTTP Proxy server running on :80")
	  http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + r.Host + r.URL.String()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
		}))	
	}()

	tlsConfig := &tls.Config{
		GetCertificate: p.getCertificate,
	}


   server := &http.Server{
	Addr: ":443",
	Handler: p.router,
	TLSConfig: tlsConfig,
   }

   log.Println("HTTPS Proxy server running on :443")
   log.Fatal(server.ListenAndServeTLS("", ""))
}



func (p *ProxyServer) loadAllCertificates() error {
	p.certCache = make(map[string]*tls.Certificate)
	for _, srv :=  range p.config {
		if srv.CertFile == "" || srv.KeyFile == "" {
			continue
		}

		cert , err := tls.LoadX509KeyPair(srv.CertFile, srv.KeyFile)
		if err != nil {
			return fmt.Errorf("TLS load failed for %s: %v", srv.Domain, err)
		}

		p.certCache[srv.Domain] = &cert
		log.Println("Loaded TLS for:", srv.Domain)
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