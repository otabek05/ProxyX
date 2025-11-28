package proxy

import (
	"ProxyX/internal/common"
	"log"
	"net/http"
)


type ProxyServer struct {
	router http.Handler
}

func NewServer(config *common.ProxyConfig) *ProxyServer {
	p := &ProxyServer{}
	p.router = NewRouter(config)
	return p
}

func (p *ProxyServer) Start()  error {
   log.Println("Listening on : 8000")
   return http.ListenAndServe(":8000", p.router)
}