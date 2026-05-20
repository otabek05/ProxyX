//go:build !windows

package cli

import (
	"ProxyX/internal/proxy"
	"log"
)

func runServer(srv *proxy.ProxyServer) {
	if err := srv.Start(); err != nil {
		log.Fatalf("ProxyX exited: %v", err)
	}
}
