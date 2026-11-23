# CLI Reference

ProxyX provides a full lifecycle management CLI for managing configuration,
TLS certificates, and service state.

---

## Available Commands

| Command | Description |
|-------|-------------|
| `config` | Show active configurations (same as `config list`) |
| `config list` | List active configurations |
| `config list -o wide` | Show full detailed configuration view |
| `config apply [file]` | Apply a configuration file |
| `config delete [name]` | Delete configuration by its **name** |
| `config example` | Print a sample server configuration |
| `config example --global` | Print the global `config.yaml` (timeouts, default rate limit, whitelist) |
| `update` | Update ProxyX to the latest GitHub release (checksum-verified, in place) |
| `cert` | Interactive TLS certificate issuance via Certbot |
| `restart` | Reload ProxyX configuration |
| `status` | Check if ProxyX service is running |
| `stop` | Stop the ProxyX service |
| `version` | Show ProxyX version |
| `healthcheck` | Configure health check endpoint (enable / disable / path / interval) |

---


### ✅ Basic CLI Usage

```bash
proxyx config example > my-app.yaml   # print a sample to start from
sudo proxyx config apply my-app.yaml
proxyx config list
proxyx config list -o wide
sudo proxyx config delete my-app
sudo proxyx restart
sudo proxyx status
```

### Wide Configuration View

```bash
proxyx config list -o wide
```
```
┌──────────────┬─────────────┬───────────┬───────────┬────────────┬──────────────┬───────────────────────┬───────────────┬─────────────────────────┐
│     FILE     │    NAME     │ NAMESPACE │  DOMAIN   │    PATH    │     TYPE     │        TARGET         │   RATELIMIT   │            TLS          │
├──────────────┼─────────────┼───────────┼───────────┼────────────┼──────────────┼───────────────────────┼───────────────┼─────────────────────────|
│ example.yaml │ local-proxy │ default   │ localhost │ /**        │ Static       │ path/to/file/         │ 1000 req /5s  │ path/to/cert/server.crt │
│              │             │           │           │            │              │                       │               │ path/to/cert/server.key │
│              │             │           │           ├────────────┼──────────────┼───────────────────────┼───────────────┼─────────────────────────┤
│              │             │           │           │ /api/v1/** │ ReverseProxy │ http://127.0.0.1:8080 │ 1000 req /5s  │ path/to/cert/server.crt │
│              │             │           │           │            │              │ http://127.0.0.1:8081 │               │ path/to/cert/server.key │
└──────────────┴─────────────┴───────────┴───────────┴────────────┴──────────────┴───────────────────────┴───────────────┴─────────────────────────┘

```
