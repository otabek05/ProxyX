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
so it behaves identically to `proxyx config ...`:

- **List** all configurations in `/etc/proxyx/conf.d/`
- **View** a configuration's YAML
- **Create / edit** a configuration in the editor (with validation)
- **Apply** — validates, writes the file, and **reloads the proxy in place**
- **Delete** — removes the file and reloads
- **Load example** — drops a ready-to-edit sample into the editor

Applying or deleting triggers an **in-process reload** — the router is rebuilt
and swapped atomically, with no process restart and no dropped connections.

## Enabling / disabling and changing the port

The panel is configured in the global config, `/etc/proxyx/config/config.yaml`:

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
| `DELETE` | `/api/configs/{name}` | Delete the configuration, reload |
| `GET` | `/api/example` | Sample server config (`?global=1` for the global config) |
