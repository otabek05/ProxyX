# ProxyX Installation Guide

This document explains how to install **ProxyX v0.1.0** on supported
Linux distributions.

------------------------------------------------------------------------

## 📦 Linux Installation

### 🐧 Debian / Ubuntu (APT-based)

#### AMD64 (x86_64)

``` bash
wget https://github.com/otabek05/ProxyX/releases/download/v0.1.2/proxyx_0.1.4_linux_amd64.deb
sudo apt install ./proxyx_0.1.2_linux_amd64.deb
```

#### ARM64

``` bash
wget https://github.com/otabek05/ProxyX/releases/download/v0.1.2/proxyx_0.1.4_linux_arm64.deb
sudo apt install ./proxyx_0.1.2_linux_arm64.deb
```

------------------------------------------------------------------------

### 🎩 RHEL / Fedora / Rocky Linux (RPM-based)

#### AMD64 (x86_64)

``` bash
wget https://github.com/otabek05/ProxyX/releases/download/v0.1.2/proxyx_0.1.4_linux_amd64.rpm
sudo dnf install ./proxyx_0.1.2_linux_amd64.rpm
```

#### ARM64

``` bash
wget https://github.com/otabek05/ProxyX/releases/download/v0.1.2/proxyx_0.1.4_linux_arm64.rpm
sudo dnf install ./proxyx_0.1.2_linux_arm64.rpm
```

------------------------------------------------------------------------

## ✅ Verify Installation

``` bash
sudo proxyx --version
```

## 🧹 Uninstall

### Debian / Ubuntu

``` bash
sudo apt remove proxyx
```

### RPM-based

``` bash
sudo dnf remove proxyx
```

------------------------------------------------------------------------

## 🍏 macOS Installation

```bash
# ---- macOS Installation ----
git clone https://github.com/otabek05/proxyx.git
cd proxyx
sudo make install-macos
sudo proxyx status
cd ..


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

------------------------------------------------------------------------

## ℹ️ Notes
-   Root privileges are required for installation
-   ProxyX binary is installed to `/usr/bin/proxyx`
-   Configuration files are located in `/etc/proxyx/`