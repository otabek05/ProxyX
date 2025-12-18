package main

import (
	"ProxyX/internal/cli"
	"ProxyX/internal/proxy"
	"ProxyX/pkg/config"
	"fmt"
	"log"
	"os"
)

func main() {
	requireRoot()
	if len(os.Args) > 1 {
		cli.Execute()
		return
	}
	
	serverConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	proxyConfig, err := config.LoadProxyXConfig()
	if err != nil {
		log.Fatalf("Failed to load proxy config: %v", err)
	}

	srv := proxy.NewServer(serverConfig, proxyConfig)
	srv.Start()
}


func requireRoot() {
	if os.Geteuid() != 0 {
		fmt.Println("This command must be run with sudo (root privileges required)")
		fmt.Println("Example: sudo proxyx services")
		os.Exit(1)
	}
}