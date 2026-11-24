// Package configmgr holds the shared logic for validating, applying, listing,
// and deleting ProxyX server configurations under conf.d. Both the CLI
// (proxyx config ...) and the web UI (proxyx web) use it, so the two frontends
// behave identically. Neither restarts the service — the caller does that.
package configmgr

import (
	"ProxyX/internal/common"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// GlobalPath is the global proxy configuration file (rate limit, whitelist,
// web panel, timeouts). It is the same path pkg/config loads at startup.
const GlobalPath = "/etc/proxyx/config/config.yaml"

// ValidateServer parses and validates a server config without writing it. It is
// used by the panel's "Check" button to report whether an uploaded or edited
// file is correct.
func ValidateServer(data []byte) (*common.ServerConfig, error) {
	var server common.ServerConfig
	if err := yaml.Unmarshal(data, &server); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}
	if err := Validate(&server); err != nil {
		return nil, err
	}
	return &server, nil
}

// RawGlobal returns the raw bytes of the global config file.
func RawGlobal() ([]byte, error) {
	return os.ReadFile(GlobalPath)
}

// ValidateGlobal parses the global config strictly (rejecting unknown keys, so
// typos are caught) without writing it.
func ValidateGlobal(data []byte) error {
	var cfg common.ProxyConfig
	if err := yaml.UnmarshalStrict(data, &cfg); err != nil {
		return fmt.Errorf("invalid global config: %w", err)
	}
	return nil
}

// ApplyGlobal validates and writes the global config file atomically. The
// service must be restarted for the new values to take effect (rate limit,
// whitelist, and web-panel settings are read only at startup), so the caller
// does not trigger an in-process reload here.
func ApplyGlobal(data []byte) error {
	if err := ValidateGlobal(data); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(GlobalPath), 0o755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	tmp := GlobalPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("failed to write global config: %w", err)
	}
	return os.Rename(tmp, GlobalPath)
}

// Summary is a lightweight view of a config for listings.
type Summary struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Domain    string `json:"domain"`
	File      string `json:"file"`
	Routes    int    `json:"routes"`
}

// Validate checks a server config's required fields and per-route shape.
func Validate(srv *common.ServerConfig) error {
	if strings.TrimSpace(srv.Metadata.Name) == "" {
		return fmt.Errorf("metadata.name is required")
	}
	if srv.Spec.Domain == "" {
		return fmt.Errorf("server missing domain")
	}
	for _, route := range srv.Spec.Routes {
		if route.Path == "" {
			return fmt.Errorf("server '%s' has route missing path", srv.Spec.Domain)
		}
		if !route.Type.IsValid() {
			return fmt.Errorf("server '%s' route '%s' has invalid type '%s'", srv.Spec.Domain, route.Path, route.Type)
		}
		switch route.Type {
		case common.RouteReverseProxy:
			if route.ReverseProxy == nil || len(route.ReverseProxy.Servers) == 0 {
				return fmt.Errorf("server '%s' route '%s' of type 'ReverseProxy' has no backends", srv.Spec.Domain, route.Path)
			}
		case common.RouteStatic:
			if route.Static == nil || route.Static.Root == "" {
				return fmt.Errorf("server '%s' route '%s' of type 'Static' missing root dir", srv.Spec.Domain, route.Path)
			}
		case common.RouteWebsocket:
			if route.Websocket == nil || route.Websocket.URL == "" {
				return fmt.Errorf("server '%s' route '%s' of type 'Websocket' missing url", srv.Spec.Domain, route.Path)
			}
		}
	}
	return nil
}

// RouteConflict reports whether newCfg collides with an existing config on the
// same domain+path. newCfgFile is the destination file for newCfg, which is
// skipped (updating a config is not a conflict with itself).
func RouteConflict(dir string, newCfg *common.ServerConfig, newCfgFile string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return err
	}
	for _, file := range files {
		if file == newCfgFile {
			continue
		}
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		var existing common.ServerConfig
		if err := yaml.Unmarshal(data, &existing); err != nil {
			return err
		}
		if existing.Spec.Domain != newCfg.Spec.Domain {
			continue
		}
		for _, newRoute := range newCfg.Spec.Routes {
			for _, oldRoute := range existing.Spec.Routes {
				if oldRoute.Path == newRoute.Path {
					return fmt.Errorf("domain='%s' path='%s' already exists in %s",
						newCfg.Spec.Domain, newRoute.Path, existing.Metadata.Name)
				}
			}
		}
	}
	return nil
}

// Apply validates data and writes it to <dir>/<metadata.name>.yaml. It returns
// the destination file path. The caller restarts the service afterwards.
func Apply(dir string, data []byte) (string, error) {
	var server common.ServerConfig
	if err := yaml.Unmarshal(data, &server); err != nil {
		return "", fmt.Errorf("invalid YAML: %w", err)
	}
	if err := Validate(&server); err != nil {
		return "", err
	}
	dest := filepath.Join(dir, server.Metadata.Name+".yaml")
	if err := RouteConflict(dir, &server, dest); err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}
	if err := os.WriteFile(dest, data, 0o644); err != nil {
		return "", fmt.Errorf("failed to write config file: %w", err)
	}
	return dest, nil
}

// Delete removes the config whose metadata.name is name. It returns the deleted
// file path, or an error if no config matched.
func Delete(dir, name string) (string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return "", err
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		var server common.ServerConfig
		if err := yaml.Unmarshal(data, &server); err != nil {
			continue
		}
		if server.Metadata.Name == name {
			if err := os.Remove(file); err != nil {
				return "", err
			}
			return file, nil
		}
	}
	return "", fmt.Errorf("configuration '%s' not found", name)
}

// List returns a summary of every config in dir.
func List(dir string) ([]Summary, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	out := make([]Summary, 0, len(files))
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		var server common.ServerConfig
		if err := yaml.Unmarshal(data, &server); err != nil {
			continue
		}
		out = append(out, Summary{
			Name:      server.Metadata.Name,
			Namespace: server.Metadata.Namespace,
			Domain:    server.Spec.Domain,
			File:      filepath.Base(file),
			Routes:    len(server.Spec.Routes),
		})
	}
	return out, nil
}

// Raw returns the YAML bytes of the config whose metadata.name is name.
func Raw(dir, name string) ([]byte, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		var server common.ServerConfig
		if err := yaml.Unmarshal(data, &server); err != nil {
			continue
		}
		if server.Metadata.Name == name {
			return data, nil
		}
	}
	return nil, fmt.Errorf("configuration '%s' not found", name)
}
