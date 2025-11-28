package main

import (
	"ProxyX/internal/cli"
	"ProxyX/internal/proxy"
	"log"
)


func main() {
	opts := cli.ParseCLI()

	if handled := cli.HandleCLI(opts); handled {
		return
	}

	config, err := proxy.LoadConfig(opts.ConfigPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := proxy.NewServer(config)
	log.Printf("Starting ProxyX on port %d ...", opts.Port)
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}

}