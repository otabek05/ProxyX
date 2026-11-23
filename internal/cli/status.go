package cli

import (
	"ProxyX/internal/pkg/config"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func (c *CLI) statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "📊 Display ProxyX runtime status",
		//PreRun: requireRoot,
		Run: func(cmd *cobra.Command, args []string) {
			c.Service.Status()
			c.printConfig()
		},
	}
}

func (c *CLI) printConfig() {

	cfg, _ := config.LoadProxyXConfig()
	fmt.Println(strings.Repeat("-", 40))
	//	fmt.Println("🛠️  ProxyX Configuration")
	//	fmt.Println(strings.Repeat("=", 40))

	// 🌐 HTTP
	fmt.Println("🌐 HTTP")
	fmt.Printf("    Read Timeout        : %s\n", cfg.HTTP.ReadTimeout)
	fmt.Printf("    Read Header Timeout : %s\n", cfg.HTTP.ReadHeaderTimeout)
	fmt.Printf("    Write Timeout       : %s\n", cfg.HTTP.WriteTimeout)
	fmt.Printf("    Idle Timeout        : %s\n", cfg.HTTP.IdleTimeout)
	fmt.Printf("    Max Header Bytes    : %d\n", cfg.HTTP.MaxHeaderBytes)

	fmt.Println(strings.Repeat("-", 40))

	// 🔒 HTTPS
	fmt.Println("🔒 HTTPS")
	fmt.Printf("    Read Timeout        : %s\n", cfg.HTTPS.ReadTimeout)
	fmt.Printf("    Read Header Timeout : %s\n", cfg.HTTPS.ReadHeaderTimeout)
	fmt.Printf("    Write Timeout       : %s\n", cfg.HTTPS.WriteTimeout)
	fmt.Printf("    Idle Timeout        : %s\n", cfg.HTTPS.IdleTimeout)
	fmt.Printf("    Max Header Bytes    : %d\n", cfg.HTTPS.MaxHeaderBytes)

	fmt.Println(strings.Repeat("-", 40))

	// ❤️ Health Check
	fmt.Println("❤️  Health Check")
	status := "Disabled"
	if cfg.HealthCheck.Enabled {
		status = "Enabled"
	}
	fmt.Printf("    Status   : %s\n", status)
	fmt.Printf("    Path     : %s\n", cfg.HealthCheck.Path)
	fmt.Printf("    Interval : %s\n", cfg.HealthCheck.Interval)

	fmt.Println(strings.Repeat("-", 40))

	// 🌐 Config Panel — on by default; off only when explicitly disabled.
	fmt.Println("🌐 Config Panel")
	if cfg.WebPanel != nil && !cfg.WebPanel.Enabled {
		fmt.Printf("    Status : Disabled\n")
	} else {
		port := 8787
		if cfg.WebPanel != nil && cfg.WebPanel.Port > 0 {
			port = cfg.WebPanel.Port
		}
		fmt.Printf("    Status : Enabled (runs with the service)\n")
		fmt.Printf("    URL    : http://127.0.0.1:%d  (localhost only)\n", port)
	}

	fmt.Println(strings.Repeat("-", 40))

	// 🔐 Certbot
	fmt.Println("🔐 Certbot")
	fmt.Printf("  Email : %s\n", cfg.Certbot.Email)
}
