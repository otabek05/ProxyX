package proxy

import (
	"net/http"
	"os"
	"path/filepath"
)

func ServerDefaultPage(w http.ResponseWriter) {
	path := filepath.Join("/etc/proxyx/web", "index.html")
	content, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "Default page not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}