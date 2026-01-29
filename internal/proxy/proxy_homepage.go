package proxy

import (
	"os"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

func ServeError(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusServiceUnavailable)
	path := filepath.Join("/etc/proxyx/web", "error.html")
	content, err := os.ReadFile(path)
	if err != nil {
		ctx.SetBodyString("No upstream")
		return 
	}

	ctx.SetContentType("text/html")
	ctx.SetBody(content)


}

func ServeDefault(ctx *fasthttp.RequestCtx) {
	path := filepath.Join("/etc/proxyx/web", "index.html")
	content, err := os.ReadFile(path)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Default page not found")
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/html")
	ctx.SetBody(content)
}




