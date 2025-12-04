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
	if len(os.Args) > 1 {
		requireRoot()
		cli.Execute()
		return
	}
	
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Println(config)

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
