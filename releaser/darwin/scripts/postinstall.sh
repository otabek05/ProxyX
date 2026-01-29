#!/bin/sh
set -e

# Create config directory
mkdir -p /usr/local/etc/proxyx

# Copy web assets
if [ -d /usr/local/etc/proxyx.web ]; then
    cp -R /usr/local/etc/proxyx.web/* /usr/local/etc/proxyx/
    rm -rf /usr/local/etc/proxyx.web
fi



# -----------------------------
if command -v proxyx >/dev/null 2>&1; then
    mkdir -p /usr/local/share/zsh/site-functions
    proxyx completion zsh > /usr/local/share/zsh/site-functions/_proxyx
fi


# Reload launchd and start service
launchctl unload /Library/LaunchDaemons/com.proxyx.service.plist 2>/dev/null || true
launchctl load /Library/LaunchDaemons/com.proxyx.service.plist
launchctl kickstart -k system/com.proxyx.service
