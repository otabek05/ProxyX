package proxy

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"
)

// acmeWebroot is where certbot (run with --webroot by the config panel) writes
// HTTP-01 challenge tokens, and where the proxy reads them from. It must match
// configmgr.AcmeWebroot.
const acmeWebroot = "/var/lib/proxyx/acme"

const acmeChallengePrefix = "/.well-known/acme-challenge/"

// maybeServeACME serves a Let's Encrypt HTTP-01 challenge on :80 if the request
// is for one, returning true if it handled the request. It must run before any
// host routing or HTTPS redirect so a domain that is mid-issuance (config
// present, no cert yet) or already on HTTPS can still be validated without
// downtime. Path traversal is blocked: the cleaned path must still live under
// the challenge prefix.
func maybeServeACME(ctx *fasthttp.RequestCtx) bool {
	raw := string(ctx.Path())
	if !strings.HasPrefix(raw, acmeChallengePrefix) {
		return false
	}
	clean := filepath.Clean(raw)
	if !strings.HasPrefix(clean, acmeChallengePrefix) {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return true
	}
	data, err := os.ReadFile(filepath.Join(acmeWebroot, clean))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return true
	}
	ctx.SetContentType("text/plain")
	ctx.SetBody(data)
	return true
}
