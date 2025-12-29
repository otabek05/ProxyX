package cli

import (
	"ProxyX/internal/common"
	"ProxyX/utils"
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func (c *CLI) certCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cert",
	    Short: "Issue TLS certificate for a domain using certbot",
		RunE: func(cmd *cobra.Command, args []string) error  {
		files, err := filepath.Glob(filepath.Join(c.serviceConfig, "*.yaml"))
		if err != nil || len(files) == 0 {
			return errors.New("No config files found.")
		}

		domainMap := make(map[int]string)
		c.printDomains(files, domainMap)
		if len(domainMap) == 0 {
			return errors.New("No domains found in configs.")
		}

		reader := bufio.NewReader(os.Stdin)
		domain, err := c.requestDomain(reader, domainMap)
		if err != nil {
			return err 
		}

		email, err := c.requestEmail(reader)
		if err != nil {
			return err 
		}

		fmt.Println("\nRequesting certificate for:", domain)
		c.Service.Stop()

		if err := c.requestCert(domain, email); err != nil {
			return fmt.Errorf("Certbot failed: %v", err )
		}

		fmt.Println("\nCertificate issued successfully!")
		c.applyCerts(&domain, files)

		fmt.Println("\nReloading ProxyX...")
		return c.Service.Restart()
		},
	}
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

func (c *CLI) requestEmail(reader *bufio.Reader) (string, error) {
	for {
		fmt.Print("Enter email for Let's Encrypt (q to exit): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("input error: %w", err)
		}

		input = strings.TrimSpace(input)

		if strings.EqualFold(input, "q") {
			return "", fmt.Errorf("user exited")
		}

		if utils.IsValidEmail(input) {
			return input, nil
		}

		fmt.Printf("Invalid email: %s. Please enter a valid email.\n", input)
		fmt.Println("Enter 'q' to exit.")
	}
}

func (c *CLI) applyCerts(domain *string, files []string) {
	certPath := "/etc/letsencrypt/live/" + *domain + "/fullchain.pem"
	keyPath := "/etc/letsencrypt/live/" + *domain + "/privkey.pem"

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}

		var server common.ServerConfig
		if err := yaml.Unmarshal(data, &server); err != nil {
			continue
		}

		if server.Spec.Domain != *domain {
			continue
		}

		server.Spec.TLS.CertFile = certPath
		server.Spec.TLS.KeyFile = keyPath
		out, _ := yaml.Marshal(&server)
		os.WriteFile(file, out, 0644)
		fmt.Println("TLS updated in:", filepath.Base(file))

	}
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
