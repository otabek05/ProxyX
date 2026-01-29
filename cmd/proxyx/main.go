package main

import (
	"ProxyX/internal/cli"
	"ProxyX/internal/platform"
	"log"
	"runtime/debug"
)

func main() {
	debug.SetGCPercent(200)
	service, err := platform.NewService()
	if err != nil {
		log.Fatal(err)
	}
	
	cmd := cli.NewCLI(service)
	cmd.Execute()
}


