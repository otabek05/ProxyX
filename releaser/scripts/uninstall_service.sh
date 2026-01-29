#!/bin/bash

sudo systemctl stop proxyx
sudo systemctl disable proxyx 
sudo rm -f /etc/systemd/system/proxyx.service
sudo rm -f /usr/local/bin/proxyx
sudo rm -rf /etc/proxyx
sudo systemctl daemon-reload

#rm -f /etc/bash_completion.d/proxyx
sudo rm -f /etc/bash_completion.d/proxyx

# Zsh
sudo rm -f /usr/local/share/zsh/site-functions/_proxyx

echo "ProxyX removed!"