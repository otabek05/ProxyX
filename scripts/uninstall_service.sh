#!/bin/bash

sudo systemctl stop proxyx
sudo systemctl disable proxyx 
sudo rm -f /etc/systemd/system/proxyx.service
sudo rm -f /usr/local/bin/proxyx
sudo rm -rf /etc/proxyx
sudo systemctl daemon-reload

echo "ProxyX removed!"