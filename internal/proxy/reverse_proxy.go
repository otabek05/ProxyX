package proxy

import (
	"crypto/tls"
	"log"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var reqPool = sync.Pool{
	New: func() interface{} { return new(fasthttp.Request) },
}

var resPool = sync.Pool{
	New: func() interface{} { return new(fasthttp.Response) },
}

func (p *ProxyServer) reverseProxyxHandler(ctx *fasthttp.RequestCtx, matched *routeInfo) {
	target := matched.loadBalancer.Next()
	if target == nil {
		ctx.Error("No upstream", fasthttp.StatusServiceUnavailable)
		return
	}

	key := target.Scheme + "://" + target.Host
	client, ok := p.proxies[key]
	if !ok {
		client = newUpstreamClient(target.Scheme == "https")
		p.proxies[key] = client
	}

	req := reqPool.Get().(*fasthttp.Request)
	defer func() {
		req.Reset()
		reqPool.Put(req)
	}()

	ctx.Request.CopyTo(req)

	uri := req.URI()
	uri.Reset()
	uri.SetScheme(target.Scheme)
	uri.SetHost(target.Host)
	uri.SetPathBytes(ctx.Path())
	uri.SetQueryStringBytes(ctx.QueryArgs().QueryString())

	resp := resPool.Get().(*fasthttp.Response)
	defer func() {
		resp.Reset()
		resPool.Put(resp)
	}()

	req.Header.Set("X-Forwarded-For", ctx.RemoteAddr().String())
	req.Header.Set("X-Forwarded-Host", string(ctx.Host()))
	req.Header.Set("X-Forwarded-Proto", map[bool]string{
		true:  "https",
		false: "http",
	}[ctx.IsTLS()])

	if err := client.DoTimeout(req, resp, 5*time.Second); err != nil {
		log.Println("upstream error:", err)
		ctx.Error("Bad Gateway", fasthttp.StatusBadGateway)
		return
	}

	resp.CopyTo(&ctx.Response)
}



func newUpstreamClient(isTLS bool) *fasthttp.Client {
	c := &fasthttp.Client{
		MaxConnsPerHost:     2048,
		MaxIdleConnDuration: 30 * time.Second,
		ReadTimeout:         15 * time.Second,
		WriteTimeout:        15 * time.Second,
	}

	if isTLS {
		c.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	c.DisableHeaderNamesNormalizing = true
	c.NoDefaultUserAgentHeader = true

	return c
}
