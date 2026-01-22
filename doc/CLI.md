# CLI Reference

ProxyX provides a full lifecycle management CLI for managing configuration,
TLS certificates, and service state.

---

## Available Commands

| Command | Description |
|-------|-------------|
| `apply` | Apply a configuration file |
| `cert` | Interactive TLS certificate issuance via Certbot |
| `config` | Show active configurations |
| `config -o wide` | Show full detailed configuration view |
| `delete` | Delete applied configuration (default behavior) |
| `delete [name]` | Delete configuration by its **name** |
| `restart` | Reload ProxyX configuration |
| `status` | Check if ProxyX service is running |
| `stop` | Stop the ProxyX service |
| `version` | Show ProxyX version |
| `healthcheck` | Configure health check endpoint (enable / disable / path / interval) |

---


### ✅ Basic CLI Usage

```bash
sudo proxyx apply  path/to/file
sudo proxyx configs
sudo proxyx configs -o wide
sudo proxyx restart
sudo proxyx status
```

### Wide Configuration View

```bash
sudo proxyx configs -o wide
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
