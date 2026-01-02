package proxy

import (
	"log"

	wsProxy "github.com/yeqown/fasthttp-reverse-proxy/v2"

	"github.com/valyala/fasthttp"
)

func (p *ProxyServer) configureWSProxy() {
	for _, c := range p.config {
		for _, r := range c.Spec.Routes {
			domain := c.Spec.Domain

			if r.Websocket == nil {
				continue
			}

			if _, exists := p.wsProxies.Load(domain); exists {
				continue
			}

		   customProxy, err := newWsProxyInstance(r.Websocket.URL)
			if err != nil {
				log.Fatal(err)
			}

			p.wsProxies.Store(domain, customProxy)
			log.Println("WebSocket proxy enabled for:", domain)

		}
	}
}

func (p *ProxyServer) websocketProxyHandler(ctx *fasthttp.RequestCtx) {

	wsServerAny, ok := p.wsProxies.Load(string(ctx.Host()))
	if !ok {
		ServeProxyHomepage(ctx)
		return
	}

	ws := wsServerAny.(*wsProxy.WSReverseProxy)

	ctx.Request.Header.Set(wsProxy.DefaultOverrideHeader, string(ctx.Path()))
	ws.ServeHTTP(ctx)
}

func newWsProxyInstance(serverURL string) (*wsProxy.WSReverseProxy, error) {
	return wsProxy.NewWSReverseProxyWith(
		wsProxy.WithURL_OptionWS(serverURL),
		wsProxy.WithDynamicPath_OptionWS(
			true,
			wsProxy.DefaultOverrideHeader,
		),
	)
}
