## 🚀 ProxyX

ProxyX is a high-performance, configuration-driven
reverse proxy and static file server written in Go.

Inspired by Nginx, ProxyX uses fasthttp instead of
Go’s standard net/http to achieve maximum performance.


## ✨ Features

- 🔁 Reverse Proxy
- 📁 Static File Serving
- 🔐 TLS / HTTPS with Certbot
- ⚖️ Load Balancing (Round-Robin)
- ❤️ Health Checking
- 🚦 Per-Domain Rate Limiting
- 📄 Declarative YAML Configuration
- 🧰 Powerful Interactive CLI Tool


## 📦 ProxyX Installation Guide

ProxyX supports the following platforms:

- 🐧 Linux
  - Debian-based (.deb)
  - RPM-based (.rpm)
- 🍎 macOS (Darwin)

📖 For detailed installation instructions:
doc/INSTALL.md


## ⚙️ ProxyX Configuration Guide

ProxyX uses a Kubernetes-style YAML configuration
to define:

- 🌍 Domains
- 🔐 TLS settings
- 🚦 Rate limits
- 🔀 Routing rules

📖 Full documentation:
doc/CONFIGURATION.md


## 🧰 ProxyX CLI Tool

ProxyX includes a full lifecycle management CLI
to manage configurations, services, and TLS.

📖 Full command reference:
doc/CLI.md


## 🔐 ProxyX TLS Configuration

ProxyX integrates with Certbot to automatically
issue and manage Let’s Encrypt TLS certificates.

📖 TLS setup guide:
doc/TLS.md


## 🖥️ System Service & Ports

ProxyX installs itself as a Linux system service
called:

- proxyx.service

It is designed to run as a production-grade daemon.


### ✅ Service Features

- 🔄 Runs as a system service
- 🚀 Starts automatically on boot
- ♻️ Restarts automatically after shutdown
- 🛑 Restarts automatically on crash or failure


### 🌐 Network Ports

- 🌍 Port 80  → HTTP traffic
- 🔐 Port 443 → HTTPS (TLS)

⚠️ Root (sudo) access is required
to bind to ports 80 and 443.


## 🧠 Architecture Overview

- ⚡ Go fasthttp server for high concurrency
- 📄 Custom YAML configuration parser
- 🔁 Reverse proxy engine
- ❤️ Health checker
- 🔐 Certbot shell integration

Middleware pipeline:
- 📝 Request Logger
- 🚦 Per-Domain Rate Limiter
- ⚖️ Load Balancer
- ❤️ Health Checker


## 🎯 Use Cases

- 🚪 API Gateway
- 🌐 Static website hosting
- 🔀 Internal microservice routing
- 🧪 Development reverse proxy
- 🏭 Production HTTPS entrypoint


## 🛡️ Security Features

- 🔐 HTTPS with Let’s Encrypt
- 🚦 Per-domain rate limiting
- ❤️ Backend health validation
- 🔒 Mandatory TLS for production


## 👤 Author

Developed by Otabek
Go Backend Developer


## 📄 License

MIT License
