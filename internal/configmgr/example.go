package configmgr

// ServerExample is a sample per-domain server configuration, printed by
// `proxyx config example` and offered by the web UI editor.
const ServerExample = `# ProxyX server configuration — save as <name>.yaml, then:
#   sudo proxyx config apply <name>.yaml
apiVersion: proxyx.io/v1
kind: ProxyConfig
metadata:
  name: my-app          # unique name; becomes <name>.yaml in /etc/proxyx/conf.d
  namespace: default
spec:
  domain: example.com

  # Optional TLS. Omit to serve plain HTTP on :80.
  tls:
    certFile: /etc/letsencrypt/live/example.com/fullchain.pem
    keyFile: /etc/letsencrypt/live/example.com/privkey.pem

  # Optional per-domain rate limit (per client IP). Omit to use the global
  # default from /etc/proxyx/config/config.yaml.
  rateLimit:
    requests: 100
    windowSeconds: 1

  routes:
    # Reverse proxy (round-robin load balancing across servers)
    - name: api
      path: /api/**
      type: ReverseProxy
      reverseProxy:
        servers:
          - url: http://127.0.0.1:8080
          - url: http://127.0.0.1:8081

    # Static file serving
    - name: assets
      path: /static/**
      type: Static
      static:
        root: /var/www/example.com

    # WebSocket passthrough
    - name: realtime
      path: /ws/**
      type: Websocket
      websocket:
        url: ws://127.0.0.1:9000
`

// GlobalExample is a sample global config.yaml, printed by
// `proxyx config example --global`.
const GlobalExample = `# ProxyX global configuration — /etc/proxyx/config/config.yaml
# Applies to the whole proxy. Created with defaults on first run; edit and
# restart ProxyX to take effect ( sudo proxyx restart ).
http:
  readTimeout: 3s
  readHeaderTimeout: 2s
  writeTimeout: 5s
  idleTimeout: 15s
  maxHeaderBytes: 1048576
https:
  readTimeout: 5s
  readHeaderTimeout: 3s
  writeTimeout: 10s
  idleTimeout: 30s
  maxHeaderBytes: 1048576
healthCheck:
  enabled: true
  path: /health
  interval: 10s
certbot:
  email: you@example.com

# Default rate limit for any domain whose own spec sets no rateLimit.
rateLimit:
  requests: 1200
  windowSeconds: 1

# Global IP whitelist.
#   mode: bypass  -> listed IPs skip rate limiting; everyone else is served
#   mode: allow   -> only listed IPs may connect; all others get 403
# Entries may be single IPs or CIDR blocks. An empty list disables the whitelist.
whitelist:
  mode: bypass
  ips:
    - 127.0.0.1
    - 10.0.0.0/8
`
