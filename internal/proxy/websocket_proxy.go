package proxy

import (
	"io"
	"log"
	"net"
	"net/url"

	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
)

func websocketProxyHandler(ctx *fasthttp.RequestCtx, matched *routeInfo) {
	u, err := url.Parse(matched.routeConfig.Websocket.URL)
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString("Invalid backend URL")
		return
	}

	var upgrader websocket.FastHTTPUpgrader
	upgrader.CheckOrigin = func(ctx *fasthttp.RequestCtx) bool { return true }

	err = upgrader.Upgrade(ctx, func(clientWs *websocket.Conn) {
		backendConn, err := net.Dial("tcp", u.Host)
		if err != nil {
			clientWs.Close()
			return
		}
		defer backendConn.Close()
		defer clientWs.Close()

		clientConn := clientWs.UnderlyingConn()
		if _, err := backendConn.Write(ctx.Request.Header.Header()); err != nil {
			return
		}

		go func() {
			_, _ = io.Copy(backendConn, clientConn)
		}()
		_, _ = io.Copy(clientConn, backendConn)
	})

	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
	}
}
