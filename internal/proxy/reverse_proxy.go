package proxy

import (
	"crypto/tls"
	"net/url"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var reqPool = sync.Pool{
	New: func() any { return new(fasthttp.Request) },
}

var resPool = sync.Pool{
	New: func() any { return new(fasthttp.Response) },
}

func (p *ProxyServer) reverseProxyxHandler(ctx *fasthttp.RequestCtx, matched *routeInfo) {
	target := matched.loadBalancer.Next()
	if target == nil {
		ctx.Error("No upstream", fasthttp.StatusServiceUnavailable)
		return
	}

	key := target.Scheme + "://" + target.Host
	clientAny, ok := p.proxies.Load(key)
	if !ok {
		client := newUpstreamClient(target)
		actual, _ := p.proxies.LoadOrStore(key, client)
		clientAny = actual
	}

	client := clientAny.(*fasthttp.HostClient)

	req := reqPool.Get().(*fasthttp.Request)
	resp := resPool.Get().(*fasthttp.Response)
	defer func() {
		req.Reset()
		reqPool.Put(req)

		resp.Reset()
		resPool.Put(resp)
	}()

	ctx.Request.CopyTo(req)

	uri := req.URI()
	uri.Reset()
	uri.SetScheme(target.Scheme)
	uri.SetHost(target.Host)
	uri.SetPathBytes(ctx.Path())
	uri.SetQueryStringBytes(ctx.QueryArgs().QueryString())

	req.Header.Set("X-Forwarded-For", ctx.RemoteAddr().String())
	req.Header.Set("X-Forwarded-Host", string(ctx.Host()))
	req.Header.Set("X-Forwarded-Proto", map[bool]string{
		true:  "https",
		false: "http",
	}[ctx.IsTLS()])

	if err := client.DoTimeout(req, resp, 5*time.Second); err != nil {
		ctx.Error("Bad Gateway", fasthttp.StatusBadGateway)
		return
	}

	resp.CopyTo(&ctx.Response)

}

func newUpstreamClient(target *url.URL) *fasthttp.HostClient {
	isTLS := target.Scheme == "https"

	c := &fasthttp.HostClient{
		Addr:                          target.Host,
		MaxConns:                      8192,
		MaxIdleConnDuration:           10 * time.Second,
		ReadTimeout:                   5 * time.Second,
		WriteTimeout:                  5 * time.Second,
		DisableHeaderNamesNormalizing: true,
		NoDefaultUserAgentHeader:      true,
	}

	if isTLS {
		c.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			ClientSessionCache:       tls.NewLRUClientSessionCache(16384),
			SessionTicketsDisabled:   false,
			PreferServerCipherSuites: true,
		}
	}

	return c
}
