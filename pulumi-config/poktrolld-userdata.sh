#!/bin/bash

# Add debug 
exec > >(tee /var/log/user-data.log|logger -t user-data -s 2>/dev/console) 2>&1

# Install Docker
sudo apt-get update -y
sudo apt-get -y install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
curl -SL https://github.com/docker/compose/releases/download/v2.26.0/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Repo setup
git clone https://github.com/pokt-network/poktroll-docker-compose-example.git /opt/poktroll
cd /opt/poktroll
curl https://raw.githubusercontent.com/pokt-network/pocket-network-genesis/master/poktrolld/testnet-validated.json > poktrolld-data/config/genesis.json
cp .env.sample .env
PUBLIC_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)
sed -i "s/YOUR_NODE_IP_OR_HOST/$PUBLIC_IP/g" .env
sed -i "s/YOUR_NODE_IP_OR_HOST/$PUBLIC_IP/g" ./relayminer-example/config/relayminer_config.yaml
sed -i "s/YOUR_NODE_IP_OR_HOST/$PUBLIC_IP/g" ./supplier_stake_config_example.yaml
/usr/local/bin/docker-compose up -d

