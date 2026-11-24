package configmgr

import (
	"ProxyX/internal/common"
	"ProxyX/internal/pkg/config"
	"ProxyX/internal/pkg/utils"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// AcmeWebroot is where certbot (run with --webroot) writes HTTP-01 challenge
// tokens and where the proxy serves them from. It must match the proxy's
// acmeWebroot constant.
const AcmeWebroot = "/var/lib/proxyx/acme"

// DomainTLS describes a domain's TLS state for the panel's TLS tab.
type DomainTLS struct {
	Domain   string `json:"domain"`
	File     string `json:"file"`
	HasTLS   bool   `json:"hasTLS"`   // the config has a tls block
	CertFile string `json:"certFile"` // configured cert path (if any)
	KeyFile  string `json:"keyFile"`  // configured key path (if any)
	CertOK   bool   `json:"certOK"`   // the cert file actually exists on disk
}

// CertEmail returns the Let's Encrypt account email from the global config, or
// "" if none is set.
func CertEmail() string {
	cfg, err := config.LoadProxyXConfig()
	if err != nil || cfg == nil {
		return ""
	}
	return cfg.Certbot.Email
}

// saveCertEmail persists the Let's Encrypt account email to the global config.
func saveCertEmail(email string) error {
	cfg, err := config.LoadProxyXConfig()
	if err != nil {
		return err
	}
	cfg.Certbot.Email = email
	return config.SaveProxyXConfig(cfg)
}

// Domains returns the TLS state of every domain across the configs in dir,
// de-duplicated by domain.
func Domains(dir string) ([]DomainTLS, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	seen := make(map[string]bool)
	out := make([]DomainTLS, 0, len(files))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		var server common.ServerConfig
		if err := yaml.Unmarshal(data, &server); err != nil {
			continue
		}
		d := server.Spec.Domain
		if d == "" || seen[d] {
			continue
		}
		seen[d] = true
		info := DomainTLS{Domain: d, File: filepath.Base(file)}
		if server.Spec.TLS != nil {
			info.HasTLS = true
			info.CertFile = server.Spec.TLS.CertFile
			info.KeyFile = server.Spec.TLS.KeyFile
			if info.CertFile != "" {
				if _, err := os.Stat(info.CertFile); err == nil {
					info.CertOK = true
				}
			}
		}
		out = append(out, info)
	}
	return out, nil
}

func domainExists(dir, domain string) bool {
	for _, d := range mustDomains(dir) {
		if d.Domain == domain {
			return true
		}
	}
	return false
}

func mustDomains(dir string) []DomainTLS {
	d, _ := Domains(dir)
	return d
}

// IssueCert obtains a Let's Encrypt certificate for domain using certbot in
// webroot mode (no downtime — the proxy serves the HTTP-01 challenge on :80),
// then writes the cert/key paths into the matching config files. The caller
// reloads the proxy afterwards so the new certificate is picked up.
func IssueCert(dir, domain, email string) error {
	if !domainExists(dir, domain) {
		return fmt.Errorf("domain %q is not in any configuration — add a server for it first", domain)
	}
	if strings.TrimSpace(email) == "" {
		email = CertEmail()
	}
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("no Let's Encrypt email set — enter one (it is saved for next time)")
	}
	if !utils.IsValidEmail(email) {
		return fmt.Errorf("invalid email address: %q", email)
	}
	_ = saveCertEmail(email)

	challengeDir := filepath.Join(AcmeWebroot, ".well-known", "acme-challenge")
	if err := os.MkdirAll(challengeDir, 0o755); err != nil {
		return fmt.Errorf("could not create ACME webroot: %w", err)
	}

	cmd := exec.Command("certbot", "certonly",
		"--webroot", "-w", AcmeWebroot,
		"-d", domain,
		"--non-interactive", "--agree-tos",
		"-m", email,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("certbot failed: %w\n%s", err, tail(out))
	}
	return writeTLSPaths(dir, domain)
}

// RenewCert renews the certificate for domain (certbot reuses the webroot
// authenticator recorded at issuance). The caller reloads the proxy afterwards.
func RenewCert(domain string) error {
	cmd := exec.Command("certbot", "renew", "--cert-name", domain, "--non-interactive")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("certbot renew failed: %w\n%s", err, tail(out))
	}
	return nil
}

// writeTLSPaths points every config for domain at the standard Let's Encrypt
// cert/key paths for that domain.
func writeTLSPaths(dir, domain string) error {
	certPath := "/etc/letsencrypt/live/" + domain + "/fullchain.pem"
	keyPath := "/etc/letsencrypt/live/" + domain + "/privkey.pem"
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return err
	}
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
		out, err := yaml.Marshal(&server)
		if err != nil {
			return err
		}
		if err := os.WriteFile(file, out, 0o644); err != nil {
			return err
		}
		updated = true
	}
	if !updated {
		return fmt.Errorf("certificate issued but no config matched domain %q", domain)
	}
	return nil
}

// tail returns the last few lines of command output for error messages.
func tail(b []byte) string {
	lines := strings.Split(strings.TrimRight(string(b), "\n"), "\n")
	if len(lines) > 8 {
		lines = lines[len(lines)-8:]
	}
	return strings.Join(lines, "\n")
}
