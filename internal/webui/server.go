// Package webui is the localhost-only browser config panel that runs inside
// the proxy service. It lets you view, create, edit, apply, and delete server
// configurations (and the global config) from a browser, plus watch live
// traffic, reusing internal/configmgr so it behaves exactly like the CLI. It
// binds to the loopback interfaces only (IPv4
// 127.0.0.1 and IPv6 ::1) and additionally refuses any non-loopback client,
// so it is never reachable from the network.
package webui

import (
	"ProxyX/internal/configmgr"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//go:embed index.html
var indexHTML []byte

const maxBodyBytes = 1 << 20 // 1 MiB is plenty for a config file

// Server serves the config panel.
type Server struct {
	Dir     string       // the conf.d directory
	Restart func() error // reloads the proxy after a write
	// Metrics returns cumulative traffic counters for the live meter. It may be
	// nil (then the meter reports zeros).
	Metrics func() (bytesIn, bytesOut, requests int64)
}

// wsUpgrader upgrades the live-metrics connection. Origin checks are left open
// because loopbackOnly already guarantees the request came from this machine.
var wsUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ListenAndServe binds to both loopback interfaces — IPv4 127.0.0.1 and IPv6
// ::1 — on the given port, so the panel is reachable whether "localhost"
// resolves to IPv4 or IPv6. It stays loopback-only. It returns when any
// listener stops, or an error if neither loopback address could be bound.
func (s *Server) ListenAndServe(port int) error {
	httpServer := &http.Server{
		Handler:           s.loopbackOnly(s.routes()),
		ReadHeaderTimeout: 5 * time.Second,
	}

	addrs := []string{
		fmt.Sprintf("127.0.0.1:%d", port),
		fmt.Sprintf("[::1]:%d", port),
	}
	var listeners []net.Listener
	var lastErr error
	for _, addr := range addrs {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			lastErr = err
			continue
		}
		listeners = append(listeners, ln)
	}
	if len(listeners) == 0 {
		return fmt.Errorf("could not bind to loopback on port %d (already in use?): %w", port, lastErr)
	}

	errc := make(chan error, len(listeners))
	for _, ln := range listeners {
		go func(l net.Listener) { errc <- httpServer.Serve(l) }(ln)
	}
	return <-errc
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.index)
	mux.HandleFunc("GET /api/configs", s.list)
	mux.HandleFunc("GET /api/configs/{name}", s.get)
	mux.HandleFunc("DELETE /api/configs/{name}", s.delete)
	mux.HandleFunc("POST /api/apply", s.apply)
	mux.HandleFunc("POST /api/validate", s.validate)
	mux.HandleFunc("GET /api/example", s.example)
	mux.HandleFunc("GET /api/global", s.getGlobal)
	mux.HandleFunc("POST /api/global", s.applyGlobal)
	mux.HandleFunc("GET /api/metrics", s.metrics)
	mux.HandleFunc("GET /api/metrics/ws", s.metricsWS)
	mux.HandleFunc("GET /api/tls", s.tlsStatus)
	mux.HandleFunc("POST /api/cert/issue", s.certIssue)
	mux.HandleFunc("POST /api/cert/renew", s.certRenew)
	return mux
}

// loopbackOnly rejects any request that did not come from the local machine.
// The listener is already bound to 127.0.0.1; this is defense in depth.
func (s *Server) loopbackOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		if ip := net.ParseIP(host); ip == nil || !ip.IsLoopback() {
			http.Error(w, "forbidden: localhost only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(indexHTML)
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	items, err := configmgr.List(s.Dir)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	data, err := configmgr.Raw(s.Dir, r.PathValue("name"))
	if err != nil {
		writeErr(w, http.StatusNotFound, err)
		return
	}
	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	_, _ = w.Write(data)
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	file, err := configmgr.Delete(s.Dir, r.PathValue("name"))
	if err != nil {
		writeErr(w, http.StatusNotFound, err)
		return
	}
	if err := s.Restart(); err != nil {
		writeErr(w, http.StatusInternalServerError, fmt.Errorf("deleted, but restart failed: %w", err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "file": file})
}

func (s *Server) apply(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes))
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	file, err := configmgr.Apply(s.Dir, body)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if err := s.Restart(); err != nil {
		writeErr(w, http.StatusInternalServerError, fmt.Errorf("saved, but restart failed: %w", err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "file": file})
}

// validate checks a config without writing it (the "Check" button). Add
// ?global=1 to validate the global config instead of a server config.
func (s *Server) validate(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes))
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if r.URL.Query().Get("global") == "1" {
		if err := configmgr.ValidateGlobal(body); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true, "kind": "global"})
		return
	}
	srv, err := configmgr.ValidateServer(body)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok": true, "kind": "server",
		"name": srv.Metadata.Name, "domain": srv.Spec.Domain, "routes": len(srv.Spec.Routes),
	})
}

// getGlobal returns the raw global config file.
func (s *Server) getGlobal(w http.ResponseWriter, r *http.Request) {
	data, err := configmgr.RawGlobal()
	if err != nil {
		writeErr(w, http.StatusNotFound, err)
		return
	}
	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	_, _ = w.Write(data)
}

// applyGlobal validates and writes the global config. It does NOT restart the
// service: rate limit, whitelist, and panel settings are read only at startup,
// so it reports that a restart is required to apply them.
func (s *Server) applyGlobal(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes))
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if err := configmgr.ApplyGlobal(body); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":              true,
		"restartRequired": true,
		"message":         "Saved. Restart the service to apply: sudo proxyx restart",
	})
}

// metrics returns the current cumulative traffic counters (poll fallback for
// the live meter when a WebSocket is unavailable).
func (s *Server) metrics(w http.ResponseWriter, r *http.Request) {
	var in, out, req int64
	if s.Metrics != nil {
		in, out, req = s.Metrics()
	}
	writeJSON(w, http.StatusOK, map[string]any{"bytesIn": in, "bytesOut": out, "requests": req})
}

// metricsWS streams traffic counters once per second over a WebSocket, adding
// per-second deltas (inRate/outRate) so the browser can show live speed.
func (s *Server) metricsWS(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Drain anything the client sends so a close is detected promptly.
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				conn.Close()
				return
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var lastIn, lastOut int64
	first := true
	for range ticker.C {
		var in, out, req int64
		if s.Metrics != nil {
			in, out, req = s.Metrics()
		}
		var inRate, outRate int64
		if !first {
			inRate = in - lastIn
			outRate = out - lastOut
		}
		first = false
		lastIn, lastOut = in, out

		_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
		if err := conn.WriteJSON(map[string]any{
			"bytesIn": in, "bytesOut": out, "requests": req,
			"inRate": inRate, "outRate": outRate,
		}); err != nil {
			return
		}
	}
}

// tlsStatus returns the current certbot email and the TLS state of every
// configured domain.
func (s *Server) tlsStatus(w http.ResponseWriter, r *http.Request) {
	domains, err := configmgr.Domains(s.Dir)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"email":   configmgr.CertEmail(),
		"webroot": configmgr.AcmeWebroot,
		"domains": domains,
	})
}

// certIssue obtains a Let's Encrypt certificate for a domain via certbot
// (webroot mode, no downtime), writes the cert paths into its config, and
// reloads the proxy to pick up the new certificate.
func (s *Server) certIssue(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domain string `json:"domain"`
		Email  string `json:"email"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, maxBodyBytes)).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		return
	}
	if err := configmgr.IssueCert(s.Dir, req.Domain, req.Email); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if err := s.Restart(); err != nil {
		writeErr(w, http.StatusInternalServerError, fmt.Errorf("certificate issued, but reload failed: %w", err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "message": "Certificate issued for " + req.Domain + " ✓"})
}

// certRenew renews a domain's certificate and reloads the proxy.
func (s *Server) certRenew(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Domain string `json:"domain"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, maxBodyBytes)).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, fmt.Errorf("invalid request: %w", err))
		return
	}
	if err := configmgr.RenewCert(req.Domain); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	if err := s.Restart(); err != nil {
		writeErr(w, http.StatusInternalServerError, fmt.Errorf("certificate renewed, but reload failed: %w", err))
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "message": "Certificate renewed for " + req.Domain + " ✓"})
}

func (s *Server) example(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/yaml; charset=utf-8")
	if r.URL.Query().Get("global") == "1" {
		_, _ = io.WriteString(w, configmgr.GlobalExample)
		return
	}
	_, _ = io.WriteString(w, configmgr.ServerExample)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
