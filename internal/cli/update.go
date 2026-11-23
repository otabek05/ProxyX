package cli

import (
	"ProxyX/internal/updater"
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func (c *CLI) updateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "🆙 Update ProxyX to the latest release",
		Long: "🆙 Download the latest ProxyX release from GitHub, verify its checksum,\n" +
			"replace this binary in place, and restart the service.\n\n" +
			"Needs root (it writes the installed binary and restarts the service).",
		Example: "  sudo proxyx update",
		PreRun:  requireRoot,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := updater.Update(version, os.Stdout)
			if err != nil {
				color.Red.Printf("❌ update failed: %s\n", err.Error())
				os.Exit(1)
			}
			if !res.Updated {
				color.Green.Printf("✅ Already up to date (v%s)\n", res.To)
				return nil
			}

			color.Green.Printf("✅ Updated %s → %s\n", res.From, res.To)
			fmt.Println("🔄 Restarting ProxyX service…")
			if err := c.Service.Restart(); err != nil {
				color.Yellow.Printf("⚠️  Binary updated, but restart failed: %s\n", err.Error())
				fmt.Println("   Run: sudo proxyx restart")
				os.Exit(1)
			}
			color.Green.Printf("🎉 ProxyX is now running v%s\n", res.To)
			return nil
		},
	}
}
