package proxy

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

// Nginx-style WebSocket proxy configuration
type WebSocketProxyConfig struct {
	// Connection settings
	ConnectTimeout    time.Duration
	SendTimeout       time.Duration
	ReadTimeout       time.Duration
	
	// Buffer sizes (nginx defaults)
	ProxyBufferSize   int // Default: 4k or 8k
	
	// Headers
	ProxySetHeaders   map[string]string
	ProxyPassHeaders  []string
	
	// Timeouts
	ProxyTimeout      time.Duration
}

func DefaultWebSocketProxyConfig() *WebSocketProxyConfig {
	return &WebSocketProxyConfig{
		ConnectTimeout:  60 * time.Second,
		SendTimeout:     60 * time.Second,
		ReadTimeout:     60 * time.Second,
		ProxyBufferSize: 8192, // 8k like nginx
		ProxyTimeout:    24 * time.Hour, // Long timeout for WebSockets
		ProxySetHeaders: map[string]string{
			"X-Real-IP":       "$remote_addr",
			"X-Forwarded-For": "$proxy_add_x_forwarded_for",
			"X-Forwarded-Proto": "$scheme",
		},
	}
}

// Main handler function (nginx-style)
func websocketProxyHandler(ctx *fasthttp.RequestCtx, matched *routeInfo) {
	config := DefaultWebSocketProxyConfig()
	
	// Verify WebSocket upgrade request
	if !isWebSocketUpgrade(ctx) {
		ctx.Error("Bad Request", fasthttp.StatusBadRequest)
		return
	}

	// Apply rate limiting
	if matched.rateLimiter != nil {
		if !matched.rateLimiter.Allow(ctx.RemoteIP().String()) {
			ctx.Error("Too Many Requests", fasthttp.StatusTooManyRequests)
			return
		}
	}

	// Get backend URL
	if matched.routeConfig.Websocket == nil || matched.routeConfig.Websocket.URL == "" {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return
	}

	backendURL := matched.routeConfig.Websocket.URL
	
	// Connect to backend with timeout
	backendConn, err := dialBackend(backendURL, config.ConnectTimeout)
	if err != nil {
		log.Printf("Backend connection failed: %v", err)
		ctx.Error("Bad Gateway", fasthttp.StatusBadGateway)
		return
	}

	// Send upgrade request to backend
	if err := sendUpgradeRequest(ctx, backendConn, backendURL, config); err != nil {
		log.Printf("Upgrade request failed: %v", err)
		backendConn.Close()
		ctx.Error("Bad Gateway", fasthttp.StatusBadGateway)
		return
	}

	// Hijack client connection
	ctx.Hijack(func(clientConn net.Conn) {
		defer clientConn.Close()
		defer backendConn.Close()

		// Receive and forward upgrade response
		if err := receiveUpgradeResponse(clientConn, backendConn, config); err != nil {
			log.Printf("Upgrade response failed: %v", err)
			return
		}

		log.Printf("WebSocket established: %s -> %s", ctx.RemoteIP(), backendURL)

		// Proxy data bidirectionally (nginx-style)
		proxyWebSocketNginxStyle(clientConn, backendConn, config)
	})
}

func dialBackend(backendURL string, timeout time.Duration) (net.Conn, error) {
	host := strings.TrimPrefix(backendURL, "ws://")
	host = strings.TrimPrefix(host, "wss://")
	usesTLS := strings.HasPrefix(backendURL, "wss://")
	
	// Extract host:port
	if idx := strings.Index(host, "/"); idx != -1 {
		host = host[:idx]
	}

	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: 30 * time.Second, // TCP keepalive like nginx
	}

	if usesTLS {
		return tls.DialWithDialer(dialer, "tcp", host, &tls.Config{
			InsecureSkipVerify: false,
		})
	}

	return dialer.Dial("tcp", host)
}

func sendUpgradeRequest(ctx *fasthttp.RequestCtx, conn net.Conn, backendURL string, config *WebSocketProxyConfig) error {
	req := &ctx.Request

	// Extract path from backend URL
	backendPath := extractPath(backendURL)
	if backendPath == "" {
		backendPath = string(req.URI().PathOriginal())
	}

	// Build request line
	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("%s %s", req.Header.Method(), backendPath))
	if len(req.URI().QueryString()) > 0 {
		buf.WriteString("?")
		buf.WriteString(string(req.URI().QueryString()))
	}
	buf.WriteString(" HTTP/1.1\r\n")

	// Nginx-style header forwarding
	req.Header.VisitAll(func(key, value []byte) {
		keyStr := string(key)
		
		// Skip Host header, we'll set it
		if strings.ToLower(keyStr) == "host" {
			return
		}
		
		buf.WriteString(keyStr)
		buf.WriteString(": ")
		buf.WriteString(string(value))
		buf.WriteString("\r\n")
	})

	// Set Host header (nginx always sets this)
	backendHost := extractHost(backendURL)
	buf.WriteString("Host: ")
	buf.WriteString(backendHost)
	buf.WriteString("\r\n")

	// Add proxy headers (nginx-style)
	for header, value := range config.ProxySetHeaders {
		resolvedValue := resolveNginxVariable(value, ctx)
		buf.WriteString(header)
		buf.WriteString(": ")
		buf.WriteString(resolvedValue)
		buf.WriteString("\r\n")
	}

	buf.WriteString("\r\n")

	// Send with timeout
	conn.SetWriteDeadline(time.Now().Add(config.SendTimeout))
	_, err := conn.Write([]byte(buf.String()))
	conn.SetWriteDeadline(time.Time{})

	return err
}

func receiveUpgradeResponse(clientConn, backendConn net.Conn, config *WebSocketProxyConfig) error {
	backendConn.SetReadDeadline(time.Now().Add(config.ReadTimeout))
	defer backendConn.SetReadDeadline(time.Time{})

	// Read response using buffered reader (nginx-style)
	reader := bufio.NewReaderSize(backendConn, config.ProxyBufferSize)
	
	// Read status line
	statusLine, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// Check for 101 Switching Protocols
	if !strings.Contains(statusLine, "101") {
		return fmt.Errorf("backend rejected upgrade: %s", statusLine)
	}

	// Read headers until empty line
	var headers []string
	headers = append(headers, statusLine)
	
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		headers = append(headers, line)
		if line == "\r\n" {
			break
		}
	}

	// Send complete response to client
	response := strings.Join(headers, "")
	_, err = clientConn.Write([]byte(response))
	if err != nil {
		return err
	}

	// Handle any buffered data
	if reader.Buffered() > 0 {
		buffered := make([]byte, reader.Buffered())
		_, err := io.ReadFull(reader, buffered)
		if err != nil {
			return err
		}
		_, err = clientConn.Write(buffered)
		if err != nil {
			return err
		}
	}

	return nil
}

func proxyWebSocketNginxStyle(clientConn, backendConn net.Conn, config *WebSocketProxyConfig) {
	var wg sync.WaitGroup
	wg.Add(2)

	// Client -> Backend
	go func() {
		defer wg.Done()
		defer func() {
			if tcp, ok := backendConn.(*net.TCPConn); ok {
				tcp.CloseWrite()
			}
		}()
		
		// Use configured buffer size (nginx default: 8k)
		buf := make([]byte, config.ProxyBufferSize)
		
		for {
			// Set read timeout (nginx proxy_read_timeout)
			clientConn.SetReadDeadline(time.Now().Add(config.ProxyTimeout))
			
			n, err := clientConn.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("Client read error: %v", err)
				}
				return
			}

			// Set write timeout (nginx proxy_send_timeout)
			backendConn.SetWriteDeadline(time.Now().Add(config.SendTimeout))
			
			_, err = backendConn.Write(buf[:n])
			if err != nil {
				log.Printf("Backend write error: %v", err)
				return
			}
		}
	}()

	// Backend -> Client
	go func() {
		defer wg.Done()
		defer func() {
			if tcp, ok := clientConn.(*net.TCPConn); ok {
				tcp.CloseWrite()
			}
		}()
		
		buf := make([]byte, config.ProxyBufferSize)
		
		for {
			backendConn.SetReadDeadline(time.Now().Add(config.ProxyTimeout))
			
			n, err := backendConn.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("Backend read error: %v", err)
				}
				return
			}

			clientConn.SetWriteDeadline(time.Now().Add(config.SendTimeout))
			
			_, err = clientConn.Write(buf[:n])
			if err != nil {
				log.Printf("Client write error: %v", err)
				return
			}
		}
	}()

	wg.Wait()
	log.Println("WebSocket connection closed gracefully")
}

// Helper functions

func isWebSocketUpgrade(ctx *fasthttp.RequestCtx) bool {
	upgrade := strings.ToLower(string(ctx.Request.Header.Peek("Upgrade")))
	connection := strings.ToLower(string(ctx.Request.Header.Peek("Connection")))
	return upgrade == "websocket" && strings.Contains(connection, "upgrade")
}

func extractPath(url string) string {
	url = strings.TrimPrefix(url, "ws://")
	url = strings.TrimPrefix(url, "wss://")
	
	if idx := strings.Index(url, "/"); idx != -1 {
		return url[idx:]
	}
	return "/"
}

func extractHost(url string) string {
	host := strings.TrimPrefix(url, "ws://")
	host = strings.TrimPrefix(host, "wss://")
	
	if idx := strings.Index(host, "/"); idx != -1 {
		return host[:idx]
	}
	return host
}

func resolveNginxVariable(value string, ctx *fasthttp.RequestCtx) string {
	// Resolve nginx-style variables
	switch value {
	case "$remote_addr":
		return ctx.RemoteIP().String()
	case "$scheme":
		if ctx.IsTLS() {
			return "https"
		}
		return "http"
	case "$proxy_add_x_forwarded_for":
		existing := string(ctx.Request.Header.Peek("X-Forwarded-For"))
		if existing != "" {
			return existing + ", " + ctx.RemoteIP().String()
		}
		return ctx.RemoteIP().String()
	default:
		return value
	}
}