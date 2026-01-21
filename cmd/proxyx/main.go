package main

import (
	"ProxyX/internal/cli"
	"ProxyX/internal/platform"
	//"ProxyX/internal/proxy"
	//"ProxyX/pkg/config"
	"fmt"
	"log"
	"os"
	"runtime/debug"
)

func main() {
	requireRoot()
	debug.SetGCPercent(200)
	service, err := platform.NewService()
	if err != nil {
		log.Fatal(err)
	}

	cmd := cli.NewCLI(service)
	if len(os.Args) == 1 {
		cmd.Help()
		return
	}

	cmd.Execute()

	/*
	
	proxyConfig, err := config.LoadProxyXConfig()
	if err != nil {
		log.Fatalf("Failed to load proxy config: %v", err)
	}

	serverConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := proxy.NewServer(serverConfig, proxyConfig)
	srv.Start()
	*/
	
}

func requireRoot() {
	if os.Geteuid() != 0 {
		fmt.Println("This command must be run with sudo (root privileges required)")
		fmt.Println("Example: sudo proxyx services")
		os.Exit(1)
	}
}
