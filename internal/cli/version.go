package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// version is injected at build time via
//
//	-ldflags "-X ProxyX/internal/cli.version=<version>"
//
// (GoReleaser sets it from the git tag). Local/manual builds show "dev".
var version = "dev"

func (c *CLI) versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "📦 Display the current ProxyX version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ProxyX version %s\n", version)
		},
	}
}
