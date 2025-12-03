#!/bin/bash

set -e 

echo "===================="
echo " Building ProxyX ..."
echo "===================="


#Build proxy server
GOOS=linux GOARCH=amd64 go build -o ../bin/proxyx ../cmd/proxyx/main.go
echo "[DONE] Built proxyx"

# Build API server
GOOS=linux GOARCH=amd64 go build -o ../bin/proxyx-api ../cmd/proxyx-api/main.go
echo "[DONE] Built proxyx-api"


if ! grep -q 'ProxyX/bin' ~/.bashrc; then
    echo 'export PATH="$HOME/ProxyX/bin:$PATH"' >> ~/.bashrc
    echo "Added ProxyX/bin to PATH in ~/.bashrc"
fi

source ~/.bashrc

sudo mkdir -p /etc/proxyx/web-admin
sudo mkdir -p /var/log/proxyx


cd ../frontend
npm install 

npm run build

echo "===================="
echo " Copying build to /etc/proxyx/web-admin ..."
echo "===================="

sudo rm -rf /etc/proxyx/web-admin/*
sudo cp -r dist/* /etc/proxyx/web-admin/

echo "[DONE] React dashboard copied to internal/web-admin"

# Install systemd services
sudo cp ../systemd/proxyx.service /etc/systemd/system/proxyx.service
sudo cp ../systemd/proxyx-api.service /etc/systemd/system/proxyx-api.service

sudo systemctl daemon-reload
sudo systemctl enable proxyx proxyx-api
sudo systemctl restart proxyx proxyx-api

echo "===================="
echo " ProxyX installed and running!"
echo " - ProxyX (proxy) : 80/443"
echo " - ProxyX API     : 5053"
echo "===================="