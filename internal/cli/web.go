package cli

import (
	"ProxyX/internal/webui"
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

// DefaultWebPort is the localhost port the config panel listens on.
const DefaultWebPort = 8787

func (c *CLI) webCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "web",
		Short: "🌐 Open the browser config panel (localhost only)",
		Long: "🌐 Start the ProxyX config panel: a small web UI, bound to 127.0.0.1,\n" +
			"for viewing, editing, applying, and deleting configurations from a browser.\n\n" +
			"It is reachable only from this machine (localhost) and never from the network.\n" +
			"Applying or deleting a config restarts the proxy, so it needs root.",
		Example: "  sudo proxyx web\n  sudo proxyx web --port 9000",
		PreRun:  requireRoot,
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")
			if port < 1 || port > 65535 {
				return fmt.Errorf("invalid port %d (must be 1-65535)", port)
			}
			srv := &webui.Server{
				Dir:     c.serviceConfig,
				Restart: c.Service.Restart,
			}
			color.Green.Printf("🌐 ProxyX config panel running — open it in a browser on THIS machine:\n")
			fmt.Printf("     http://127.0.0.1:%d   (or http://localhost:%d)\n", port, port)
			fmt.Println("   Localhost only. If you're on a remote server, tunnel first:")
			fmt.Printf("     ssh -L %d:127.0.0.1:%d user@server\n", port, port)
			fmt.Println("   Press Ctrl+C to stop.")
			if err := srv.ListenAndServe(port); err != nil {
				color.Red.Printf("❌ config panel stopped: %s\n", err.Error())
				os.Exit(1)
			}
			return nil
		},
	}
	cmd.Flags().IntP("port", "p", DefaultWebPort, "Localhost port for the config panel")
	return cmd
}
