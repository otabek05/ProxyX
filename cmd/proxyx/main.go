package main

import (
	"ProxyX/internal/cli"
	"ProxyX/internal/platform"
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
	cmd.Execute()
}

func requireRoot() {
	if os.Geteuid() != 0 {
		fmt.Println("This command must be run with sudo (root privileges required)")
		os.Exit(1)
	}
}
