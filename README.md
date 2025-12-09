# ProxyX

ProxyX is a high‑performance, configuration‑driven reverse proxy and static file server written in **Go**, inspired by **Nginx**. It provides:

* ✅ Reverse Proxy
* ✅ Static File Serving
* ✅ TLS/HTTPS with Certbot
* ✅ Load Balancing (Round‑Robin)
* ✅ Health Checking
* ✅ **Per‑Domain Rate Limiting**
* ✅ Declarative YAML Configuration
* ✅ Powerful Interactive CLI Tool

---

## Features Overview

### Reverse Proxy

Route traffic to one or more backend servers with automatic load balancing and health checks.

### Static File Hosting

Serve static files directly from any directory on your system.

### TLS / HTTPS (Certbot)

Automatically secure domains using Let's Encrypt via **Certbot**.

### Load Balancing

* **Round‑Robin** distribution
* **Health Checking** for backend servers
* Automatic failover

### Per‑Domain Rate Limiting

Each domain has its **own independent rate limit**.

### YAML Configuration

Kubernetes‑style declarative configuration format.

---

## Example ProxyX Configuration

```yaml
apiVersion: proxyx.io/v1
kind: ProxyConfig

metadata:
  name: local-proxy
  namespace: default

spec:
  domain: localhost

  tls:
    certFile: /path/to/certs/server.crt
    keyFile:  /path/to/certs/server.key

  rateLimit:
    requests: 1000
    windowSeconds: 2

  routes:
    - name: api-route
      path: /api/v1/**
      type: ReverseProxy
      reverseProxy:
        servers:
          - url: http://localhost:8080
          - url: http://localhost:8081

    - name: static-files
      path: /**
      type: Static
      static:
        root: /path/to/static/file
```

---

## Route Types

### ✅ Static Route

```yaml
type: Static
static:
  root: /var/www/app
```

* Direct disk file serving
* Supports recursive path matching using `/**`

---

### ✅ Reverse Proxy Route

```yaml
type: ReverseProxy
reverseProxy:
  servers:
    - url: http://localhost:8080
    - url: http://localhost:8081
```

* Multiple backends supported
* Round‑Robin load balancing
* Automatic health‑based failover

---

## Load Balancer

### ✅ Round‑Robin

Distributes requests evenly across all **healthy** backends.

### ✅ Health Checker

* Removes offline servers automatically
* Periodic TCP/HTTP availability probing

---

## 🚦 Per‑Domain Rate Limiter

Each domain controls its **own request limits**:

```yaml
rateLimit:
  requests: 1000
  windowSeconds: 2
```

* Protects domains independently
* Prevents cross‑domain poisoning
* Applied across **all routes under the domain**

---

## TLS & HTTPS with Certbot

ProxyX integrates with **Certbot** to automatically issue and manage Let's Encrypt TLS certificates.

### ✅ Requirements

You **must install Certbot manually**:

```bash
sudo dnf install certbot   # RHEL / Amazon Linux
sudo apt install certbot   # Ubuntu / Debian
```

---

### ✅ Interactive Certificate Issuance

```bash
sudo proxyx certs
```

ProxyX will **prompt interactively**:

* ✅ Domain name
* ✅ Email address for Let's Encrypt

Then ProxyX will:

* Request the certificate
* Store the cert & key
* Automatically wire it into your configuration

---

## CLI Tool

ProxyX includes a full lifecycle management CLI.

### ✅ Available Commands

| Command           | Description                                     |
| ----------------- | ----------------------------------------------- |
| `apply`           | Apply configuration file                        |
| `certs`           | **Interactive TLS issuance via Certbot**        |
| `configs`         | Show active configurations                      |
| `configs -o wide` | Show full detailed configuration                |
| `delete`          | Delete applied configuration (default behavior) |
| `delete [name]`   | Delete configuration by its **name**            |
| `restart`         | Reload ProxyX configuration                     |
| `status`          | Check if ProxyX is running                      |
| `stop`            | Stop ProxyX service                             |
| `version`         | Show ProxyX version                             |


---

### ✅ Basic CLI Usage

```bash
sudo proxyx apply -f path/to/file
sudo proxyx configs
sudo proxyx configs -o wide
sudo proxyx restart
sudo proxyx status
```

---

## Wide Configuration View Example

```bash
sudo proxyx configs -o wide
```

```
┌──────────────┬─────────────┬───────────┬───────────┬────────────┬──────────────┬───────────────────────┬───────────────┬─────────────────────────┐
│     FILE     │    NAME     │ NAMESPACE │  DOMAIN   │    PATH    │     TYPE     │        TARGET         │   RATELIMIT   │            TLS          │
├──────────────┼─────────────┼───────────┼───────────┼────────────┼──────────────┼───────────────────────┼───────────────┼─────────────────────────|
│ example.yaml │ local-proxy │ default   │ localhost │ /**        │ Static       │     path/to/file/     │ 1000 req /5s  │ path/to/cert/server.crt │
│              │             │           │           │            │              │                       │               │ path/to/cert/server.key │
│              │             │           │           ├────────────┼──────────────┼───────────────────────┼───────────────┼─────────────────────────┤
│              │             │           │           │ /api/v1/** │ ReverseProxy │ http://localhost:8080 │ 1000 req / 5s │ path/to/cert/server.crt │
│              │             │           │           │            │              │ http://localhost:8081 │               │ path/to/cert/server.key │
└──────────────┴─────────────┴───────────┴───────────┴────────────┴──────────────┴───────────────────────┴───────────────┴─────────────────────────┘
```

---

## System Service & Ports

ProxyX automatically installs itself as a **Linux system service (`proxyx.service`)** and is designed to run as a **production-grade daemon**.

### ✅ Service Features

* ✅ Runs as `proxyx` system service
* ✅ Automatically starts on system boot
* ✅ Automatically restarts if the server turns off/on
* ✅ Automatically restarts on crash or failure

### ✅ Network Ports

* ✅ **Port 80** → HTTP traffic
* ✅ **Port 443** → HTTPS (TLS via Certbot)

> ⚠️ ProxyX requires **root (sudo)** access to bind to ports **80 and 443**.

---

## Architecture Overview

* Go `net/http` server
* Custom YAML parser
* Reverse proxy engine
* Health checker
* Certbot shell integration
* Middleware pipeline:

  * Request Logger
  * Per‑Domain Rate Limiter
  * Load Balancer
  * Health Checker

---

## Use Cases

* API Gateway
* Static website hosting
* Internal microservice router
* Development reverse proxy
* Production HTTPS entrypoint


---

## Installation

ProxyX supports **Linux** and **macOS**.

```bash
# ---- Linux Installation ----
git clone https://github.com/otabek05/proxyx.git
cd proxyx
sudo make install
sudo proxyx status
cd ..

# ---- macOS Installation ----
git clone https://github.com/otabek05/proxyx.git
cd proxyx
sudo make install-macos
sudo proxyx status
cd ..

```

---

## Uninstallation

To completely remove ProxyX from your system:

## Linux

```bash
# ---- Manual Removal ----
sudo systemctl stop proxyx
sudo systemctl disable proxyx

sudo rm -f /etc/systemd/system/proxyx.service
sudo rm -f /usr/local/bin/proxyx
sudo rm -rf /etc/proxyx
sudo systemctl daemon-reload

# ---- Removal via Makefile ----
cd ~/proxyx
sudo make uninstall
```

---

## MacOS

```bash
# ---- Manual Removal ----
  sudo launchctl stop proxyx
  sudo launchctl unload /Library/LaunchDaemons/proxyx.plist

  sudo rm -f /Library/LaunchDaemons/proxyx.plist
  sudo rm -f /usr/local/bin/proxyx
  sudo rm -rf /etc/proxyx

# ---- Removal via Makefile ----
  cd ~/proxyx 
  sudo make uninstall-macos
```

✅ ProxyX is now fully removed from your system.

---


## Security Features

* HTTPS with Let's Encrypt
* Per‑domain rate limiting
* Backend health validation
* Mandatory TLS for production

---

## Author

Developed by **Otabek** — Go Backend Developer

---

## 📄 License

MIT License
