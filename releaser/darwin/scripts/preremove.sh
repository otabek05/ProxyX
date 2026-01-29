#!/bin/sh
set -e

launchctl bootout system/com.proxyx.service 2>/dev/null || true

rm -f /usr/local/share/zsh/site-functions/_proxyx
rm -f /Library/LaunchDaemons/com.proxyx.service.plist
rm -f /usr/local/bin/proxyx
rm -rf /usr/local/etc/proxyx
