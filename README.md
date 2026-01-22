# ProxyX

ProxyX is a high-performance, configuration-driven reverse proxy and
static file server written in Go, inspired by Nginx.

It uses fasthttp instead of Go’s standard net/http to deliver
high concurrency and low latency. 🚀


## ✨ Features

- Reverse Proxy
- Static File Serving
- TLS / HTTPS with Certbot
- Load Balancing (Round-Robin)
- Health Checking
- Per-Domain Rate Limiting
- Declarative YAML Configuration
- Interactive CLI Tool


## 📦 Installation Guide

ProxyX supports the following platforms:

- Linux
  - Debian-based (.deb)
  - RPM-based (.rpm)
- macOS (Darwin)

For detailed installation steps, see:
Installation Guide → doc/INSTALL.md


## ⚙️ Configuration Guide

ProxyX uses a Kubernetes-style YAML configuration format to define:

- Domains
- TLS settings
- Rate limits
- Routing rules

For full configuration examples and explanations, see:
Configuration Guide → doc/CONFIGURATION.md


## 🧰 CLI Tool

ProxyX includes a full lifecycle management CLI for:

- Applying and deleting configs
- Managing TLS certificates
- Controlling the ProxyX service

📖 Full command reference:
CLI Tool Guide → doc/CLI.md


## 🔐 TLS Configuration

ProxyX integrates with Certbot to automatically issue and manage
Let’s Encrypt TLS certificates.

- HTTPS for Static, ReverseProxy, and WebSocket routes
- Multiple domains supported
- Automatic certificate renewal

📖 For detailed instructions, see:
TLS Configuration Guide → doc/TLS.md


## 🖥️ System Service & Ports

ProxyX installs itself as a system service:

  proxyx.service

It is designed to run as a production-grade daemon.


### ✅ Service Features

- Starts automatically on system boot
- Restarts automatically after shutdown
- Restarts on crash or failure


### 🌐 Network Ports

- Port 80  → HTTP
- Port 443 → HTTPS (TLS)

⚠️ Root (sudo) access is required to bind to ports 80 and 443.


## 🧠 Architecture Overview

- fasthttp server for high-performance HTTP handling
- Custom YAML configuration parser
- Reverse proxy engine
- Health checking subsystem
- Certbot shell integration

Middleware pipeline:
- Request Logger
- Per-Domain Rate Limiter
- Load Balancer
- Health Checker


## 🎯 Use Cases

- API Gateway
- Static website hosting
- Internal microservice router
- Development reverse proxy
- Production HTTPS entrypoint


## 🔒 Security Features

- HTTPS with Let’s Encrypt
- Per-domain rate limiting
- Backend health validation
- TLS-first production design


## 👨‍💻 Author

Developed by **Otabek**  
Go Backend Developer


## 📄 License

MIT License
