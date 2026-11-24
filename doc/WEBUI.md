# Web Config Panel

ProxyX serves a small browser-based config panel **as part of the proxy
service** — it runs in the background whenever ProxyX is running. There is no
separate command to start; just start the service.

```bash
sudo proxyx start      # the panel comes up with the proxy
```

Then open it in a browser **on the same machine**:

```
http://127.0.0.1:8787      (or http://localhost:8787)
```

`proxyx status` shows the panel's URL and whether it is enabled.

## What it does

The panel is a thin UI over the same logic the CLI uses (`internal/configmgr`),
so it behaves identically to `proxyx config ...`. It has three tabs.

### 🗂 Servers

Per-domain configurations in `/etc/proxyx/conf.d/`:

- **List** all configurations
- **View** a configuration's YAML
- **Create / edit** a configuration in the editor
- **Upload file** — load a `.yaml` from disk straight into the editor
- **Check** — validate without saving (reports name, domain, route count, or
  the exact error)
- **Apply** — validates, writes the file, and **reloads the proxy in place**
- **Delete** — removes the file and reloads
- **Load example** — drops a ready-to-edit sample into the editor

Applying or deleting triggers an **in-process reload** — the router is rebuilt
and swapped atomically, with no process restart and no dropped connections.

### ⚙️ Global

Edit the global config (`/etc/proxyx/config/config.yaml`) — rate limit, IP
whitelist, web-panel settings, timeouts — with the same Upload / Check / Load
example buttons. Saving validates **strictly** (unknown keys are rejected, so
typos are caught) and writes the file. Unlike server configs, global settings
are read only at startup, so the panel tells you to restart the service
(`sudo proxyx restart`) for the change to take effect.

### 🔐 TLS

Issue and renew Let's Encrypt certificates with certbot, straight from the
browser. The table lists every configured domain with its TLS state (no cert /
configured / cert present on disk) and an **Issue** / **Re-issue** and **Renew**
button each. Set the Let's Encrypt account email once (it is saved to the global
config for next time).

A **Challenge** selector picks how ownership is proven:

- **HTTP-01 (webroot, default)** — the proxy answers the challenge on port 80
  itself (served before any routing or HTTPS redirect), so there is **no
  downtime**. Needs the domain pointing at this server and port 80 reachable
  from the internet.
- **DNS-01** — a certbot **DNS plugin** publishes a `TXT` record at
  `_acme-challenge.<domain>`. Needs **no inbound ports** (works behind NAT or a
  firewall) and is the only method that can issue **wildcard** certificates. Give
  the plugin name (e.g. `dns-cloudflare`, `dns-route53`, `dns-google`) and the
  absolute path to its credentials file. The plugin name is validated against
  `^dns-[a-z0-9-]+$`, and certbot is run without a shell, so the value can never
  be a command. Some plugins (e.g. `dns-route53` with an IAM role) need no
  credentials file — leave it blank.

On success the standard `/etc/letsencrypt/live/<domain>/` cert and key paths are
written into the domain's config and the proxy reloads in place. Renewal reuses
whichever challenge method was recorded at issuance.

> Requirements: `certbot` installed (plus the DNS plugin package for DNS-01,
> e.g. `python3-certbot-dns-cloudflare`). For HTTP-01, DNS pointing here and port
> 80 reachable. Any certbot error (rate limit, DNS, unreachable, bad token) is
> shown in full.

This complements the CLI `proxyx cert <domain>`, which instead stops the service
and uses certbot `--standalone`.

### 🎨 Theme

A 🌙 / ☀️ toggle in the header switches between dark and light. The choice is
remembered per browser; the initial theme follows your OS preference.

### 📈 Traffic

A live, server-wide traffic meter: incoming and outgoing **speed** (per second),
cumulative totals, request count, and a rolling 60-second chart. Press **Start**
to open a WebSocket stream (`/api/metrics/ws`) that pushes a sample every second.
Figures are approximate — request/response headers plus any in-memory body —
which is accurate enough to watch load in real time. Counters are cumulative for
the life of the proxy process and survive in-place reloads.

## Enabling / disabling and changing the port

The panel is **on by default** — including when the `webPanel` block is absent
(for example in a `config.yaml` written by an older version). It is turned off
only by explicitly setting `enabled: false`.

It is configured in the global config, `/etc/proxyx/config/config.yaml`:

```yaml
webPanel:
  enabled: true   # set false to turn the panel off
  port: 8787      # localhost port
```

Print a sample with `proxyx config example --global`. After changing these
values, restart the service (`sudo proxyx restart`).

## Security

- **Localhost only.** The panel binds to the loopback interfaces only — IPv4
  `127.0.0.1` and IPv6 `::1` — and additionally rejects any request whose client
  address is not a loopback address (`403`). It is never reachable from the
  network, only from the machine itself.
- **Runs as root**, because it lives in the proxy service (which binds :80/:443
  and writes `/etc/proxyx/conf.d/`). Anyone with a local shell on the machine can
  reach it — the same trust level as running the CLI with sudo.

> To reach the panel from another machine, do **not** expose the port. Use an
> SSH tunnel: `ssh -L 8787:127.0.0.1:8787 user@server`, then open
> `http://127.0.0.1:8787` locally.

## HTTP API

The UI is served from `/`; it talks to this JSON API (all localhost only):

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/configs` | List configurations (name, domain, file, route count) |
| `GET` | `/api/configs/{name}` | Raw YAML of one configuration |
| `POST` | `/api/apply` | Body is YAML; validate, write, reload |
| `POST` | `/api/validate` | Validate a config without writing (`?global=1` for the global config) |
| `DELETE` | `/api/configs/{name}` | Delete the configuration, reload |
| `GET` | `/api/example` | Sample server config (`?global=1` for the global config) |
| `GET` | `/api/global` | Raw YAML of the global config |
| `POST` | `/api/global` | Validate (strict) and write the global config; restart required to apply |
| `GET` | `/api/metrics` | Current cumulative traffic counters (poll fallback) |
| `GET` | `/api/metrics/ws` | WebSocket stream of traffic counters + per-second rates |
| `GET` | `/api/tls` | Certbot email and per-domain TLS status |
| `POST` | `/api/cert/issue` | `{domain,email,method,plugin,creds}` — issue a cert (HTTP-01 webroot or DNS-01 plugin), write paths, reload |
| `POST` | `/api/cert/renew` | `{domain}` — renew a cert and reload |
