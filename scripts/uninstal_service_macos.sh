#!/bin/bash

echo "Uninstalling ProxyX from macOS..."

sudo launchctl stop org.proxyx.service
sudo launchctl unload /Library/LaunchDaemons/org.proxyx.service.plist

sudo rm -f /Library/LaunchDaemons/org.proxyx.service.plist
sudo rm -f /usr/local/bin/proxyx
sudo rm -rf /etc/proxyx

echo "ProxyX macOS service removed!"
