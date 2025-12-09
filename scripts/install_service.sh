#!/bin/bash
set -e

sudo cp ./bin/linux/proxyx /usr/local/bin/proxyx
sudo chmod +x /usr/local/bin/proxyx

sudo mkdir -p /etc/proxyx
sudo cp -r web /etc/proxyx

sudo cp systemd/proxyx.service /etc/systemd/system/proxyx.service

sudo systemctl daemon-reload
sudo systemctl enable proxyx
sudo systemctl restart proxyx
