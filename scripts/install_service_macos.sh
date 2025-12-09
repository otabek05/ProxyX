#!/bin/bash
set -e

sudo cp ./bin/macos/proxyx /usr/local/bin/proxyx
sudo chmod +x /usr/local/bin/proxyx

sudo mkdir -p /etc/proxyx
sudo cp -r web /etc/proxyx

sudo cp systemd/proxyx.plist /Library/LaunchDaemons/org.proxyx.service.plist
sudo chmod 644 /Library/LaunchDaemons/org.proxyx.service.plist

sudo launchctl load /Library/LaunchDaemons/org.proxyx.service.plist
sudo launchctl start org.proxyx.service
