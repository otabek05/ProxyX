package cli

import (
	"ProxyX/internal/common"
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"ProxyX/internal/pkg/config"
	"ProxyX/internal/pkg/utils"


	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func (c *CLI) certCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cert [domain | renew]",
		Short: "🔐 Manage TLS certificates with Certbot",
		Args:  cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.runCert(args); err != nil {
				fmt.Println(err)
			}
		},
	}
}



func (c *CLI) runCert( args []string) error {
	files, err := filepath.Glob(filepath.Join(c.serviceConfig, "*.yaml"))
	if err != nil || len(files) == 0 {
		return errors.New("⚠️ no configuration files found")
	}

	if len(args) == 0 {
		return c.runCertInteractive(files)
	}

	if args[0] == "renew" {
		if len(args) != 2 {
			return errors.New("❌ usage: proxyx cert renew example.com")
		}
		return c.runRenew(args[1], files)
	}

	return c.runIssueCert(args[0], files)
}


func (c *CLI) domainExists(domain string, files []string) bool {
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var server common.ServerConfig
		if yaml.Unmarshal(data, &server) != nil {
			continue
		}

		if server.Spec.Domain == domain {
			return true
		}
	}
	return false
}

func (c *CLI) runIssueCert(domain string, files []string) error {
	if !c.domainExists(domain, files) {
		return fmt.Errorf("❌ domain not found in configs: %s", domain)
	}

	reader := bufio.NewReader(os.Stdin)

	email, err := c.getCertbotEmail(reader)
	if err != nil {
		return err
	}

	fmt.Println("🔐 Requesting certificate for:", domain)
	c.Service.Stop()

	if err := c.requestCert(domain, email); err != nil {
		return fmt.Errorf("❌ certbot failed: %w", err)
	}

	fmt.Println("✅ Certificate issued successfully")
	c.applyCerts(domain, files)

	fmt.Println("🔄 Reloading ProxyX...")
	return c.Service.Restart()
}

func (c *CLI) runRenew(domain string, files []string) error {
	if !c.domainExists(domain, files) {
		return fmt.Errorf("❌ domain not found in configs: %s", domain)
	}

	fmt.Println("♻️ Renewing certificate for:", domain)

	cmd := exec.Command("certbot", "renew", "--cert-name", domain)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("✅ Certificate renewed successfully")
	c.applyCerts(domain, files)

	fmt.Println("🔄 Reloading ProxyX...")
	return c.Service.Restart()
}

func (c *CLI) runCertInteractive(files []string) error {
	domainMap := make(map[int]string)
	c.printDomains(files, domainMap)

	if len(domainMap) == 0 {
		return errors.New("⚠️ no domains found in configs")
	}

	reader := bufio.NewReader(os.Stdin)
	domain, err := c.requestDomain(reader, domainMap)
	if err != nil {
		return err
	}

	return c.runIssueCert(domain, files)
}

func (c *CLI) getCertbotEmail(reader *bufio.Reader) (string, error) {
	cfg, err := config.LoadProxyXConfig()
	if err != nil {
		return "", err
	}

	if cfg.Certbot.Email != "" {
		fmt.Println("📧 Using saved email:", cfg.Certbot.Email)
		return cfg.Certbot.Email, nil
	}

	fmt.Print("Enter email for Let's Encrypt: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	if !utils.IsValidEmail(email) {
		return "", errors.New("❌ invalid email address")
	}

	cfg.Certbot.Email = email
	_ = config.SaveProxyXConfig(cfg)

	fmt.Println("✅ Email saved for future certificate requests")
	return email, nil
}

func (c *CLI) requestDomain(reader *bufio.Reader, domainMap map[int]string) (string, error) {
	for {
		fmt.Print("\nSelect domain number (q to exit): ")
		choiceStr, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("input error: %w", err)
		}

		choiceStr = strings.TrimSpace(choiceStr)
		if strings.EqualFold(choiceStr, "q") {
			return "", fmt.Errorf("user exited")
		}

		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		domain, exists := domainMap[choice]
		if !exists {
			fmt.Println("Invalid selection.")
			continue
		}

		return domain, nil
	}
}

func (c *CLI) applyCerts(domain string, files []string) error {
	certPath := "/etc/letsencrypt/live/" + domain + "/fullchain.pem"
	keyPath := "/etc/letsencrypt/live/" + domain + "/privkey.pem"

	updated := false
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var server common.ServerConfig
		if err := yaml.Unmarshal(data, &server); err != nil {
			continue
		}

		if server.Spec.Domain != domain {
			continue
		}

		if server.Spec.TLS == nil {
			server.Spec.TLS = &common.TLSConfig{}
		}

		server.Spec.TLS.CertFile = certPath
		server.Spec.TLS.KeyFile = keyPath

		out, _ := yaml.Marshal(&server)
		os.WriteFile(file, out, 0644)
		fmt.Println("✅ TLS updated in:", filepath.Base(file))
		updated = true

	}

	if !updated {
		return fmt.Errorf("no config matched domain: %s", domain)
	}

	return nil
}

func (c *CLI) requestCert(domain, email string) error {
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
	return certCmd.Run()
}

func (c *CLI) printDomains(files []string, domainMap map[int]string) {
	fmt.Println("\nAvailable Domains:")
	fmt.Println("-------------------------")

	seen := make(map[string]bool)
	index := 1

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			fmt.Println("Failed to read file:", file, err)
			continue
		}

		var server common.ServerConfig
		if err := yaml.Unmarshal(data, &server); err != nil {
			fmt.Println("Invalid YAML:", file, err)
			continue
		}

		if _, ok := seen[server.Spec.Domain]; ok {
			continue
		}

		domainMap[index] = server.Spec.Domain
		seen[server.Spec.Domain] = true
		fmt.Printf("[%d] %s\n", index, server.Spec.Domain)
		index++
	}
}
