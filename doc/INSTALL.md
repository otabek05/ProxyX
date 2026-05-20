# ProxyX Installation Guide

This document explains how to install **ProxyX** on supported platforms.
Replace `VERSION` in the commands below with the release you want
(e.g. `0.1.4`).

------------------------------------------------------------------------

## 📦 Linux Installation

### 🐧 Debian / Ubuntu (APT-based)

#### AMD64 (x86_64)

``` bash
wget https://github.com/otabek05/ProxyX/releases/download/vVERSION/proxyx_VERSION_linux_amd64.deb
sudo apt install ./proxyx_VERSION_linux_amd64.deb
```

#### ARM64

``` bash
wget https://github.com/otabek05/ProxyX/releases/download/vVERSION/proxyx_VERSION_linux_arm64.deb
sudo apt install ./proxyx_VERSION_linux_arm64.deb
```

------------------------------------------------------------------------

### 🎩 RHEL / Fedora / Rocky Linux (RPM-based)

#### AMD64 (x86_64)

``` bash
wget https://github.com/otabek05/ProxyX/releases/download/vVERSION/proxyx_VERSION_linux_amd64.rpm
sudo dnf install ./proxyx_VERSION_linux_amd64.rpm
```

#### ARM64

``` bash
wget https://github.com/otabek05/ProxyX/releases/download/vVERSION/proxyx_VERSION_linux_arm64.rpm
sudo dnf install ./proxyx_VERSION_linux_arm64.rpm
```

#### ✅ Verify Installation

``` bash
sudo proxyx --version
```

#### 🧹 Uninstall

``` bash
sudo apt remove proxyx     # Debian / Ubuntu
sudo dnf remove proxyx     # RHEL / Fedora
```

------------------------------------------------------------------------

## 🍏 macOS Installation

### From release archive (recommended)

Pick the archive for your CPU — `arm64` for Apple Silicon (M1/M2/M3),
`amd64` for Intel Macs.

``` bash
# Apple Silicon
curl -LO https://github.com/otabek05/ProxyX/releases/download/vVERSION/proxyx_VERSION_darwin_arm64.tar.gz
tar -xzf proxyx_VERSION_darwin_arm64.tar.gz

# Install binary and config/web assets
sudo install -m 0755 proxyx /usr/local/bin/proxyx
sudo mkdir -p /etc/proxyx
sudo cp -R web /etc/proxyx/

# Optional: register as a launchd service
sudo proxyx status     # verify CLI works
```

To run ProxyX as a managed service on macOS, clone the repo and use
the bundled Makefile target (it installs the launchd plist):

``` bash
git clone https://github.com/otabek05/ProxyX.git
cd ProxyX
sudo make install-macos
sudo proxyx status
```

### 🧹 Uninstall (macOS)

``` bash
# If installed via Makefile
cd ProxyX && sudo make uninstall-macos

# Manual cleanup
sudo launchctl stop proxyx
sudo launchctl unload /Library/LaunchDaemons/proxyx.plist
sudo rm -f /Library/LaunchDaemons/proxyx.plist
sudo rm -f /usr/local/bin/proxyx
sudo rm -rf /etc/proxyx
```

------------------------------------------------------------------------

## 🪟 Windows Installation

ProxyX runs as a native Windows service via the Service Control
Manager (SCM). Pick `amd64` for x64 PCs or `arm64` for ARM devices
(Surface Pro X, Windows-on-ARM laptops).

### Install (PowerShell, **as Administrator**)

``` powershell
# 1. Download and unpack the release archive
Invoke-WebRequest `
  -Uri https://github.com/otabek05/ProxyX/releases/download/vVERSION/proxyx_VERSION_windows_amd64.zip `
  -OutFile proxyx.zip
Expand-Archive .\proxyx.zip -DestinationPath .\proxyx

# 2. Run the bundled installer
cd .\proxyx
.\scripts\install.ps1
```

The installer:

- Copies `proxyx.exe` to `C:\Program Files\ProxyX\`
- Copies `web\` assets to `C:\ProgramData\ProxyX\web\`
- Creates configuration directory `C:\ProgramData\ProxyX\conf.d\`
- Registers the `proxyx` Windows service (auto-start, restart-on-failure)
- Starts the service immediately

### Manage the service

``` powershell
proxyx status      # query SCM + show PID
proxyx restart     # stop + start via SCM
proxyx stop        # stop the service
sc.exe start proxyx   # native equivalent
```

`status`, `stop`, and `restart` all talk to the Windows Service
Control Manager and require an elevated shell.

### 🧹 Uninstall (Windows)

From the unpacked release directory in an elevated PowerShell:

``` powershell
.\scripts\uninstall.ps1
```

This stops and removes the service and deletes
`C:\Program Files\ProxyX\`. Config under `C:\ProgramData\ProxyX\`
is left intact — delete it manually if no longer needed.

------------------------------------------------------------------------

## ℹ️ Notes

- Root / Administrator privileges are required for installation and
  for any command that touches the service (`start`, `stop`,
  `restart`, `apply`, etc.).
- ProxyX binds to ports **80** and **443** by default.
- Configuration files live in:
  - Linux: `/etc/proxyx/conf.d/`
  - macOS: `/etc/proxyx/conf.d/`
  - Windows: `C:\ProgramData\ProxyX\conf.d\`
