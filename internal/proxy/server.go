package proxy

import (
	"ProxyX/internal/common"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)


type ProxyServer struct {
	router http.Handler
	config *common.ProxyConfig
}

func NewServer(config *common.ProxyConfig) *ProxyServer {
	p := &ProxyServer{config: config,}
	p.router = NewRouter(config)
	return p
}

func (p *ProxyServer) Start()  {
   certManager := autocert.Manager{
	Prompt: autocert.AcceptTOS,
	Cache: autocert.DirCache("certs"),
	HostPolicy: autocert.HostWhitelist(
	   p.config.ToDomainList()...
	),
   }

   go func ()  {
	 http.ListenAndServe(":80", certManager.HTTPHandler(nil))
   }()
   

   server := &http.Server{
	Addr: ":443",
	Handler: p.router,
	TLSConfig: certManager.TLSConfig(),
   }

   log.Println("🚀 HTTPS Proxy server running on :443")
   err := server.ListenAndServeTLS("", "")
   if err != nil {
		log.Fatal(err)
   }
}