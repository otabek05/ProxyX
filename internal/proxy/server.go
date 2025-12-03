package proxy

import (
	"ProxyX/internal/common"
	"log"
	"net/http"
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
   server := &http.Server{
	Addr: ":8000",
	Handler: p.router,
   }

   log.Println("🚀 HTTPS Proxy server running on :8000")
   err := server.ListenAndServe()
   if err != nil {
		log.Fatal(err)
   }
}