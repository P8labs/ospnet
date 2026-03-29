#!/usr/bin/env bash
set -euo pipefail

AGENT_USER="root"
AGENT_GROUP="root"
INSTALL_BIN_PATH="/usr/local/bin/ospnet-agent"
SERVICE_PATH="/etc/systemd/system/ospnet-agent.service"
CONFIG_DIR="/etc/ospnet"
DATA_DIR="/var/lib/ospnet"
ENV_FILE="/etc/default/ospnet-agent"

AGENT_VERSION="${OSPNET_AGENT_VERSION:-latest}"
AGENT_DOWNLOAD_URL="${OSPNET_AGENT_DOWNLOAD_URL:-https://ospnet.run/releases/${AGENT_VERSION}/ospnet-agent-linux-amd64}"
MASTER_URL="${OSPNET_MASTER_URL:-}"
ONBOARD_TOKEN="${OSPNET_ONBOARD_TOKEN:-}"

if [[ "$(id -u)" -ne 0 ]]; then
  echo "This installer must run as root."
  exit 1
fi

if ! command -v curl >/dev/null 2>&1; then
  apt-get update
  apt-get install -y curl ca-certificates
fi

if ! command -v docker >/dev/null 2>&1; then
  curl -fsSL https://get.docker.com | sh
fi

if ! command -v tailscale >/dev/null 2>&1; then
  curl -fsSL https://tailscale.com/install.sh | sh
fi

mkdir -p "${CONFIG_DIR}" "${DATA_DIR}"
chmod 700 "${CONFIG_DIR}"
chmod 755 "${DATA_DIR}"

if [[ -z "${MASTER_URL}" ]]; then
  read -r -p "Enter OSPNet master URL (e.g. http://100.64.0.10:8080): " MASTER_URL
fi

if [[ -z "${ONBOARD_TOKEN}" ]]; then
  read -r -p "Enter onboarding token: " ONBOARD_TOKEN
fi

cat >"${ENV_FILE}" <<EOF
OSPNET_MASTER_URL=${MASTER_URL}
OSPNET_CONFIG_PATH=/etc/ospnet/config.json
OSPNET_TOKEN_PATH=/etc/ospnet/token
OSPNET_DB_PATH=/var/lib/ospnet/agent.db
OSPNET_AGENT_PORT=9000
EOF

echo -n "${ONBOARD_TOKEN}" >"${CONFIG_DIR}/token"
chmod 600 "${CONFIG_DIR}/token"

curl -fsSL "${AGENT_DOWNLOAD_URL}" -o "${INSTALL_BIN_PATH}"
chmod +x "${INSTALL_BIN_PATH}"

cat >"${SERVICE_PATH}" <<'EOF'
[Unit]
Description=OSPNet Node Agent
After=network-online.target docker.service tailscaled.service
Wants=network-online.target docker.service tailscaled.service

[Service]
Type=simple
User=root
Group=root
EnvironmentFile=-/etc/default/ospnet-agent
ExecStart=/usr/local/bin/ospnet-agent
Restart=always
RestartSec=5
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable --now ospnet-agent

echo "OSPNet agent installed and started."
