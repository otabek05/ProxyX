package main

import (
	"ProxyX/internal/cli"
	"ProxyX/internal/proxy"
	"log"
	"os"
)



func main() {
	if len(os.Args) > 1 {
		cli.Execute()
		return
	}
	
	config, err := proxy.LoadConfig(ConfigPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := proxy.NewServer(config)
	srv.Start(); 
}