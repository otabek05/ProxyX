# Web Config Panel

ProxyX ships a small browser-based config panel so you can manage configurations
without hand-editing YAML and running CLI commands.

```bash
sudo proxyx web                 # http://127.0.0.1:8787
sudo proxyx web --port 9000     # choose another localhost port
```

Open the printed URL in a browser on the same machine.

## What it does

The panel is a thin UI over the same logic the CLI uses (`internal/configmgr`),
so it behaves identically to `proxyx config ...`:

- **List** all configurations in `/etc/proxyx/conf.d/`
- **View** a configuration's YAML
- **Create / edit** a configuration in the editor (with validation)
- **Apply** — validates, writes the file, and restarts ProxyX
- **Delete** — removes the file and restarts ProxyX
- **Load example** — drops a ready-to-edit sample into the editor

Applying or deleting restarts the proxy, exactly like the CLI.

## Security

- **Localhost only.** The server binds to the loopback interfaces only — IPv4
  `127.0.0.1` and IPv6 `::1` — and additionally rejects any request whose client
  address is not a loopback address (`403`). It is never reachable from the
  network, only from the machine itself. Open `http://127.0.0.1:8787` (or
  `http://localhost:8787`) in a browser on the same machine.
- **Runs on demand.** The panel runs only while `proxyx web` is in the
  foreground. Stop it with `Ctrl+C`. It is not part of the always-on proxy
  service and opens no port on `0.0.0.0`.
- **Requires root.** Because applying/deleting writes to `/etc/proxyx/conf.d/`
  and restarts the service, the command needs `sudo`.

> To reach the panel from another machine, do **not** expose this port. Use an
> SSH tunnel instead, e.g.
> `ssh -L 8787:127.0.0.1:8787 user@server`, then open `http://127.0.0.1:8787`
> locally.

## HTTP API

The UI is served from `/`; it talks to this JSON API (all localhost only):

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/configs` | List configurations (name, domain, file, route count) |
| `GET` | `/api/configs/{name}` | Raw YAML of one configuration |
| `POST` | `/api/apply` | Body is YAML; validate, write, restart |
| `DELETE` | `/api/configs/{name}` | Delete the configuration, restart |
| `GET` | `/api/example` | Sample server config (`?global=1` for the global config) |
