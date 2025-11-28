#!/bin/bash

set -e 

echo "Installing ProxyX..."

sudo cp ./bin/proxyx /usr/local/bin/proxyx
sudo chmod +x /usr/local/bin/proxyx


# create config file 
sudo mkdir -p /etc/proxyx
sudo cp configs/proxy.yaml /etc/proxyx/proxy.yaml

#install service 
sudo cp proxyx.service /etc/systemd/system/proxyx.service

#releaod and enable 
sudo systemctl daemon-reload
sudo systemctl enable proxyx
sudo systemctl restart proxyx

echo "Proxyx installed and running"
