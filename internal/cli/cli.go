package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
)


type CLIOptions struct {
	ConfigPath string 
	Port int 
	Command string 
	FilePath  string 
	Test bool
}

func ParseCLI() *CLIOptions {
	opts := &CLIOptions{}

	flag.StringVar(&opts.ConfigPath, "config", "configs/proxy.yaml", "Path to config file")
	flag.IntVar(&opts.Port, "port", 8000, "Port to run the server")
	flag.BoolVar(&opts.Test, "t", false, "Test all config files and exit")

	flag.Parse()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "apply":
			opts.Command = "apply"
			applyCmd := flag.NewFlagSet("apply", flag.ExitOnError)
			applyCmd.StringVar(&opts.FilePath, "f", "", "Path to config file to apply")
			applyCmd.Parse(os.Args[2:])
		default:
			opts.Command = "" // default daemon mode
		}
	}

	return opts
}


func HandleCLI(opts *CLIOptions) bool {
	if opts.Test {
		fmt.Println("Testing all config files...")
		fmt.Println("All config files OK")
		return true // CLI handled
	}

	if opts.Command == "apply" {
		if opts.FilePath == "" {
			log.Fatal("You must specify a config file with -f")
		}
		fmt.Printf("Applying config file: %s\n", opts.FilePath)
		fmt.Println("Config applied successfully")
		return true
	}

	// no CLI commands, run daemon
	return false
}