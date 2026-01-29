#!/bin/bash
set -e

sudo cp ./dist/bin/linux/proxyx /usr/local/bin/proxyx
sudo chmod +x /usr/local/bin/proxyx

sudo mkdir -p /etc/proxyx
sudo cp -r web /etc/proxyx

sudo cp releaser/linux/systemd/proxyx.service /etc/systemd/system/proxyx.service


sudo systemctl daemon-reload
sudo systemctl enable proxyx
sudo systemctl restart proxyx


# Bash completion
if command -v bash >/dev/null; then
  sudo mkdir -p /etc/bash_completion.d
  sudo proxyx completion bash | sudo tee /etc/bash_completion.d/proxyx > /dev/null
  source ~/.bashrc
fi

# Zsh completion
if command -v zsh >/dev/null; then
  sudo mkdir -p /usr/local/share/zsh/site-functions
  sudo proxyx completion zsh | sudo tee /usr/local/share/zsh/site-functions/_proxyx > /dev/null
fi