## 🔐 TLS & HTTPS with Certbot

ProxyX integrates with Certbot to issue and manage
Let’s Encrypt TLS certificates automatically.

HTTPS is supported for:
- 📁 Static routes
- 🔁 Reverse Proxy routes
- 🔌 WebSocket routes


## ✅ Requirements

Before enabling TLS, Certbot must be installed.

📦 RHEL / CentOS / Amazon Linux
sudo dnf install certbot

📦 Ubuntu / Debian
sudo apt install certbot

⚠️ Make sure:
- 🌐 The server has public internet access
- 🌍 The domain DNS record points to this server


## 📜 Certificate Issuance

ProxyX offers multiple ways to manage certificates.


### 🧭 Interactive Mode (Recommended)

sudo proxyx cert

You will be prompted to:
- 🔢 Select a domain from existing configurations
- 📧 Enter an email address for Let’s Encrypt


### 🌐 Issue Certificate for a Domain

sudo proxyx cert example.com

Notes:
- ✔️ Domain must exist in ProxyX configuration
- 🔁 Stored email will be reused automatically


### ♻️ Renew Certificate

sudo proxyx cert renew example.com

Notes:
- 🔄 Uses Certbot renewal flow
- ⚙️ Configuration is updated automatically


## ⚙️ Automatic Actions by ProxyX

After issuing or renewing a certificate, ProxyX will:

- 🔐 Request or renew the TLS certificate
- 📝 Update the configuration with:
  - certFile
  - keyFile
- 🔄 Restart the ProxyX service

🚫 No manual HTTPS configuration required.


## 📌 Notes & Best Practices

- 🔓 Ports 80 (HTTP) and 443 (HTTPS) must be open
- 📍 Use stable, public domains to avoid validation errors
- 🧩 Each domain has its own certificate
- 🌍 Certificates are trusted by all major browsers
- 💾 Email address is saved and reused automatically
