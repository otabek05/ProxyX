package cli

import (
	"ProxyX/internal/common"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)


func init(){
	rootCmd.AddCommand(certCmd)
}


var certCmd = &cobra.Command{
	Use: "cert",
	Short: "Issue TLS certificate for a domain using certbot",
	Run: func(cmd *cobra.Command, args []string) {
		configDir := "/etc/proxyx/configs"
		files, err := filepath.Glob(filepath.Join(configDir, "*.yaml"))
		if err != nil || len(files) == 0 {
			fmt.Println("No config files found")
			return
		}

		domainMap := make(map[int]string)
		index := 1

		fmt.Println("\nAvailable Domains:")
		fmt.Println("-------------------------")

		for _, file := range files {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			var cfg common.ProxyConfig
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				continue
			}

			seen := make(map[string]bool)
			for _, srv := range cfg.Servers {
				if seen[srv.Domain] {
					continue
				}

				domainMap[index] = srv.Domain
				seen[srv.Domain] = true
				fmt.Printf("[%d] %s\n", index, srv.Domain)
				index++
			}
		}

		if len(domainMap) == 0 {
			fmt.Println("No domains found in configs.")
			return
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nSelect domain number: ")
		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)

		choice, err := strconv.Atoi(choiceStr)
		if err != nil || domainMap[choice] == "" {
			fmt.Println("Invalid selection.")
			return
		}

		domain := domainMap[choice]

		fmt.Print("Enter email for Let's Encrypt: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		fmt.Println("\nRequesting certificate for:", domain)
		stopProxyX()
		certCmd := exec.Command(
			"certbot", "certonly",
			"--standalone",
			"-d", domain,
			"--non-interactive",
			"--agree-tos",
			"-m", email,
		)

		certCmd.Stdout = os.Stdout
		certCmd.Stderr = os.Stderr

		if err := certCmd.Run(); err != nil {
			fmt.Println("Certbot failed:", err)
			return
		}

		fmt.Println("\nCertificate issued successfully!")
		certPath := "/etc/letsencrypt/live/" + domain + "/fullchain.pem"
		keyPath := "/etc/letsencrypt/live/" + domain + "/privkey.pem"

		for _, file := range files {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			var cfg common.ProxyConfig
			if err := yaml.Unmarshal(data, &cfg); err != nil {
				continue
			}

			updated := false

			for i := range cfg.Servers {
				if cfg.Servers[i].Domain == domain {
					cfg.Servers[i].CertFile = certPath
					cfg.Servers[i].KeyFile = keyPath
					updated = true
				}
			}

			if updated {
				out, _ := yaml.Marshal(&cfg)
				os.WriteFile(file, out, 0644)
				fmt.Println("TLS updated in:", filepath.Base(file))
			}
		}

		fmt.Println("\nReloading ProxyX...")
		reloadProxyX()
	},
}
