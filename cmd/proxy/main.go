package main

import (
	"ProxyX/configs"
	"ProxyX/internal/cli"
	"ProxyX/internal/proxy"
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
	
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := proxy.NewServer(config)
	srv.Start()
}


func requireRoot() {
	if os.Geteuid() != 0 {
		fmt.Println("This command must be run with sudo (root privileges required)")
		fmt.Println("Example: sudo proxyx services")
		os.Exit(1)
	}
}
