package proxy

import (
	"os"
	"path/filepath"
	"github.com/valyala/fasthttp"
)

func staticRouteHandler(ctx *fasthttp.RequestCtx, matched *routeInfo) {
	if matched.routeConfig.Static == nil || matched.routeConfig.Static.Root == "" {
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetBodyString("Proxy Homepage")
		return
	}

	staticDir := filepath.Join(matched.routeConfig.Static.Root)
	requestedFile := filepath.Join(staticDir, string(ctx.Path()))
	if info, err := os.Stat(requestedFile); err == nil && !info.IsDir() {
		fasthttp.ServeFile(ctx, requestedFile)
		return
	}

	indexFile := filepath.Join(staticDir, "index.html")
	if _, err := os.Stat(indexFile); os.IsNotExist(err) {
		ServeProxyHomepage(ctx)
		return
	}

	fasthttp.ServeFile(ctx, indexFile)
}