package proxy

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/valyala/fasthttp"
)

func staticRouteHandler(ctx *fasthttp.RequestCtx, matched *routeInfo) {
	if matched.routeConfig.Static == nil || matched.routeConfig.Static.Root == "" {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBodyString("Proxy Homepage")
		return
	}

	staticDir, err := filepath.Abs(matched.routeConfig.Static.Root)
	if err != nil {
		ServeError(ctx)
		return
	}

	requestedFile, ok := safeJoin(staticDir, string(ctx.Path()))
	if !ok {
		ServeError(ctx)
		return
	}

	if info, err := os.Stat(requestedFile); err == nil && !info.IsDir() {
		fasthttp.ServeFile(ctx, requestedFile)
		return
	}

	indexFile := filepath.Join(staticDir, "index.html")
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		ServeError(ctx)
		return
	}

	fasthttp.ServeFile(ctx, indexFile)
}

// safeJoin joins reqPath onto root and ensures the result stays within root.
func safeJoin(root, reqPath string) (string, bool) {
	joined := filepath.Join(root, filepath.Clean("/"+reqPath))
	rel, err := filepath.Rel(root, joined)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", false
	}
	return joined, true
}