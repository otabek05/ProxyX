// Package webui is the localhost-only browser config panel started by
// `proxyx web`. It lets you view, create, edit, apply, and delete server
// configurations from a browser, reusing internal/configmgr so it behaves
// exactly like the CLI. It binds to 127.0.0.1 and additionally refuses any
// non-loopback client, so it is never reachable from the network.
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
)

//go:embed index.html
var indexHTML []byte

const maxBodyBytes = 1 << 20 // 1 MiB is plenty for a config file

// Server serves the config panel.
type Server struct {
	Dir     string       // the conf.d directory
	Restart func() error // restarts the proxy after a write
}

// ListenAndServe binds to 127.0.0.1:port and serves until an error occurs.
func (s *Server) ListenAndServe(port int) error {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           s.loopbackOnly(s.routes()),
		ReadHeaderTimeout: 5 * time.Second,
	}
	return httpServer.ListenAndServe()
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.index)
	mux.HandleFunc("GET /api/configs", s.list)
	mux.HandleFunc("GET /api/configs/{name}", s.get)
	mux.HandleFunc("DELETE /api/configs/{name}", s.delete)
	mux.HandleFunc("POST /api/apply", s.apply)
	mux.HandleFunc("GET /api/example", s.example)
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
