# ProxyX

ProxyX is a lightweight and flexible **reverse proxy** written in **Golang**, designed for clean, structured YAML-based configuration.  
It includes a powerful CLI tool for managing configs, validating formats, and controlling the ProxyX systemd service.

Future improvements include **automatic TLS via Let’s Encrypt** and a **web-based dashboard**.

---

## ✨ Features

- Simple YAML-based configuration
- Reverse proxy with load balancing
- Static file serving
- CLI tool (`proxyx`)
- Systemd service integration
- Strict configuration validation
- Auto cleanup (removes YAML comments)
- **Coming Soon:**  
  - Let's Encrypt TLS automation  
  - Web management UI  
  - Hot reload without restart  

---

## 📦 Installation

```bash
git clone https://github.com/otabek05/proxyx.git
cd proxyx
make install
```

# ProxyX CLI Usage

## 📘 CLI Usage

ProxyX provides a command-line interface for managing the proxy service, validating configuration files, and checking system status.

### **Basic Format**
```bash
proxyx <command> [options]
```

---

## 🔧 Commands

### **Start ProxyX**
```bash
proxyx version
```

---

### **Stop ProxyX service**
```bash
proxyx stop
```
Stops the systemd service: `proxyx.service`.

---

### **Restart ProxyX service**
```bash
proxyx restart
```
Restarts the service cleanly.

---

### **Check ProxyX status**
```bash
proxyx status
```
Shows whether ProxyX is running and displays real-time process information:

Example output:
```
ProxyX is running (systemd service)
PID       CPU%    MEM%    Uptime
3993856   0.3     1.1     01:22:55
```

### Apply a configuration file

```
proxyx apply -f /path/to/config.yaml

```
Applies the configuration file to ProxyX and validates it before saving to /etc/proxyx/configs/.
Reloads the running ProxyX service if active.


---


# ProxyX YAML Configuration Guide

This document explains the full YAML configuration format used by **ProxyX**.

---

## Top-Level Structure

A ProxyX configuration file must contain a list of servers:

```yaml
servers:
  - domain: example.com
    routes:
      - path: /
        type: proxy
        backends:
          - http://localhost:8080
```

---

## Server Structure

Each server block maps an incoming domain to one or more route definitions:

```yaml
servers:
  - domain: yourdomain.com
    routes:
      - path: /api
        type: proxy
        backends:
          - http://localhost:9000
      - path: /static
        type: static
        dir: /var/www/html
```


### **Fields**

| Field  | Required | Type | Description |
|--------|----------|------|-------------|
| `domain` | Yes | string | The hostname that this rule applies to. |
| `routes` | Yes | list | A list of routing rules. |


---

## Route Structure

Each route defines how ProxyX handles requests.

### **Common Fields**

| Field | Required | Type | Description |
|-------|----------|------|-------------|
| `path` | Yes | string | URL path prefix (e.g., `/api`). |
| `type` | Yes | string | Either `proxy` or `static`. |

---

## Proxy Route Format

Proxy routes forward requests to backend servers.

```yaml
- path: /api
  type: proxy
  backends:
    - http://127.0.0.1:4000
    - http://127.0.0.1:5000
```

### **Rules**

- `backends` **must NOT be empty**.
- Load balancing is automatically round-robin.

---

## Static Route Format

Static routes serve files from a directory.

```yaml
- path: /
  type: static
  dir: /var/proxyx/site
```

### **Rules**

- `dir` **must NOT be empty**.
- Must point to an existing directory.

---

Applying Configuration via CLI

ProxyX provides a CLI command to apply a YAML configuration and validate its format before applying:

proxyx apply -f /path/to/config.yaml

Behavior:

1. Reads the YAML file.
2. Validates structure:
   - servers must not be empty
   - Each server must have a domain
   - Routes must have valid type, path
   - proxy routes require at least 1 backend
   - static routes require a non-empty dir
3. Saves the validated file into /etc/proxyx/configs/
4. Reloads the running ProxyX service if active

Example output:

Configuration applied successfully: /etc/proxyx/configs/config.yaml


## Full Example (Valid YAML)

```yaml
servers:
  - domain: example.com
    routes:
      - path: /api
        type: proxy
        backends:
          - http://localhost:8080
      - path: /
        type: static
        dir: /usr/share/nginx/html
```

---

## Invalid Examples

### Missing backends when type = proxy

```yaml
type: proxy
backends: []
```

### Missing dir when type = static

```yaml
type: static
dir:
```

---

## ✔️ Formatting Rules

- Must follow YAML indentation (2 spaces recommended)
- No trailing comments allowed
- No trailing slashes unless intended
- Paths must begin with `/`