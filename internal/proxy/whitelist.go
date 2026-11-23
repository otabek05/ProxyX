package proxy

import (
	"net"

	"ProxyX/internal/common"
)

// Whitelist is a parsed, request-time view of common.WhitelistConfig. It holds
// individual IPs and CIDR networks so matching never re-parses strings.
type Whitelist struct {
	mode string
	ips  map[string]struct{}
	nets []*net.IPNet
}

// NewWhitelist builds a matcher from config. A nil config, or one with no mode,
// defaults to "bypass" with an empty list (i.e. a no-op that blocks nobody).
func NewWhitelist(cfg *common.WhitelistConfig) *Whitelist {
	w := &Whitelist{
		mode: common.WhitelistBypass,
		ips:  make(map[string]struct{}),
	}
	if cfg == nil {
		return w
	}
	if cfg.Mode == common.WhitelistAllow {
		w.mode = common.WhitelistAllow
	}
	for _, entry := range cfg.IPs {
		if _, network, err := net.ParseCIDR(entry); err == nil {
			w.nets = append(w.nets, network)
			continue
		}
		if ip := net.ParseIP(entry); ip != nil {
			w.ips[ip.String()] = struct{}{}
		}
	}
	return w
}

// Mode reports the configured whitelist mode ("bypass" or "allow").
func (w *Whitelist) Mode() string { return w.mode }

// Empty reports whether the whitelist has no entries. An empty "allow" list
// would block everyone, so callers treat it as disabled instead.
func (w *Whitelist) Empty() bool { return len(w.ips) == 0 && len(w.nets) == 0 }

// Contains reports whether ip (a bare address, no port) is listed.
func (w *Whitelist) Contains(ip string) bool {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return false
	}
	if _, ok := w.ips[parsed.String()]; ok {
		return true
	}
	for _, network := range w.nets {
		if network.Contains(parsed) {
			return true
		}
	}
	return false
}
