#!/usr/bin/env bash
set -euo pipefail

# ─────────────────────────────────────────────────
# Nebula VPN — Node Provisioning Script
# Usage: ./provision-node.sh <server-ip> [ssh-port] [ssh-user]
# ─────────────────────────────────────────────────

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'
log()  { echo -e "${GREEN}[✓]${NC} $1"; }
warn() { echo -e "${YELLOW}[!]${NC} $1"; }
err()  { echo -e "${RED}[✗]${NC} $1"; exit 1; }
info() { echo -e "${CYAN}[i]${NC} $1"; }

SERVER_IP="${1:-${NODE_IP:-}}"
SSH_PORT="${2:-${NODE_SSH_PORT:-22}}"
SSH_USER="${3:-${NODE_SSH_USER:-root}}"
NODE_PASSWORD="${NODE_PASSWORD:-}"

COUNTRY="${NODE_COUNTRY:-FI}"
CITY="${NODE_CITY:-Helsinki}"
LATITUDE="${NODE_LATITUDE:-60.1699}"
LONGITUDE="${NODE_LONGITUDE:-24.9384}"
BW_MBPS="${NODE_BANDWIDTH:-1000}"
MAX_CONNS="${NODE_MAX_CONNS:-1000}"

API_URL="${API_URL:-https://api.nebula-vpn.io}"
NODE_TOKEN="${NODE_TOKEN:-}"

SSH_OPTS="-o StrictHostKeyChecking=accept-new -o ConnectTimeout=10 -p $SSH_PORT"
SSH_CMD="ssh $SSH_OPTS $SSH_USER@$SERVER_IP"
SCP_CMD="scp $SSH_OPTS"

check_prereqs() {
    for cmd in ssh ssh-keygen curl; do
        command -v "$cmd" >/dev/null 2>&1 || err "Missing: $cmd"
    done
}

setup_ssh_key() {
    local key_file="$HOME/.ssh/nebula_provision"
    
    if [ ! -f "$key_file" ]; then
        info "Generating SSH key for provisioning..."
        ssh-keygen -t ed25519 -f "$key_file" -N "" -C "nebula-provision" >/dev/null 2>&1
    fi
    
    if [ -n "$NODE_PASSWORD" ]; then
        info "Copying SSH key to server (password auth)..."
        command -v sshpass >/dev/null 2>&1 || brew install sshpass >/dev/null 2>&1 || true
        if command -v sshpass >/dev/null 2>&1; then
            sshpass -p "$NODE_PASSWORD" ssh-copy-id -o StrictHostKeyChecking=accept-new \
                -p "$SSH_PORT" -i "$key_file.pub" "$SSH_USER@$SERVER_IP" >/dev/null 2>&1 || \
                warn "ssh-copy-id failed, will use password for each command"
        else
            warn "sshpass not available, installing..."
            brew install sshpass >/dev/null 2>&1
            sshpass -p "$NODE_PASSWORD" ssh-copy-id -o StrictHostKeyChecking=accept-new \
                -p "$SSH_PORT" -i "$key_file.pub" "$SSH_USER@$SERVER_IP" >/dev/null 2>&1 || true
        fi
    fi
    
    SSH_CMD="ssh -i $key_file $SSH_OPTS $SSH_USER@$SERVER_IP"
    SCP_CMD="scp -i $key_file $SSH_OPTS"
}

check_os() {
    info "Checking OS..."
    local os
    os=$($SSH_CMD "source /etc/os-release 2>/dev/null && echo \$ID || echo unknown")
    log "OS: $os"
}

install_docker() {
    info "Installing Docker..."
    $SSH_CMD "command -v docker >/dev/null 2>&1" && { log "Docker already installed"; return 0; }
    
    $SSH_CMD bash -s <<'DOCKER'
        set -e
        apt-get update -qq
        apt-get install -y -qq ca-certificates curl
        install -m 0755 -d /etc/apt/keyrings
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
        chmod a+r /etc/apt/keyrings/docker.asc
        echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
        apt-get update -qq
        apt-get install -y -qq docker-ce docker-ce-cli containerd.io docker-compose-plugin
        systemctl enable docker
        systemctl start docker
DOCKER
    log "Docker installed"
}

install_wireguard() {
    info "Installing WireGuard..."
    $SSH_CMD "command -v wg >/dev/null 2>&1 && command -v wg-quick >/dev/null 2>&1" && { log "WireGuard already installed"; return 0; }
    
    $SSH_CMD bash -s <<'WG'
        set -e
        apt-get update -qq
        apt-get install -y -qq wireguard-tools linux-headers-$(uname -r) 2>/dev/null || \
        apt-get install -y -qq wireguard-tools
        sysctl -w net.ipv4.ip_forward=1
        sysctl -w net.ipv6.conf.all.forwarding=1
        echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf
        echo "net.ipv6.conf.all.forwarding=1" >> /etc/sysctl.conf
WG
    log "WireGuard installed"
}

generate_wireguard_keys() {
    info "Generating WireGuard keys..."
    local keys
    keys=$($SSH_CMD "wg genkey | tee /etc/nebula/wg_private.key | wg pubkey | tee /etc/nebula/wg_public.key" 2>&1) || true
    
    SERVER_PUBKEY=$($SSH_CMD "cat /etc/nebula/wg_public.key 2>/dev/null")
    if [ -z "$SERVER_PUBKEY" ]; then
        $SSH_CMD "mkdir -p /etc/nebula && wg genkey | tee /etc/nebula/wg_private.key | wg pubkey | tee /etc/nebula/wg_public.key" >/dev/null 2>&1
        SERVER_PUBKEY=$($SSH_CMD "cat /etc/nebula/wg_public.key")
    fi
    log "WireGuard public key: $SERVER_PUBKEY"
}

configure_firewall() {
    info "Configuring firewall..."
    $SSH_CMD bash -s <<'FW'
        set -e
        ufw --force reset >/dev/null 2>&1
        ufw default deny incoming
        ufw default allow outgoing
        ufw allow ssh
        ufw allow 51820/udp comment 'WireGuard'
        ufw allow 9100/tcp comment 'node-exporter'
        ufw allow 9586/tcp comment 'wireguard-exporter'
        ufw --force enable >/dev/null 2>&1
FW
    log "Firewall configured"
}

setup_node_config() {
    info "Setting up node configuration..."
    $SSH_CMD bash -s -- "$COUNTRY" "$CITY" "$LATITUDE" "$LONGITUDE" "$BW_MBPS" "$MAX_CONNS" <<'CONFIG'
        set -e
        COUNTRY="$1"; CITY="$2"; LATITUDE="$3"; LONGITUDE="$4"; BW="$5"; MAXC="$6"
        mkdir -p /opt/nebula/config /opt/nebula/data /var/log/nebula
        
        cat > /opt/nebula/config/node.json <<EOF
{
  "server_id": "$(hostname)-$(date +%s)",
  "host": "$(curl -4 -s ifconfig.me 2>/dev/null || curl -4 -s icanhazip.com 2>/dev/null || echo 'unknown')",
  "ipv4": "$(curl -4 -s ifconfig.me 2>/dev/null || echo 'unknown')",
  "port": 51820,
  "protocol": "wireguard",
  "public_key": "$(cat /etc/nebula/wg_public.key 2>/dev/null || echo 'unknown')",
  "country": "$COUNTRY",
  "city": "$CITY",
  "latitude": $LATITUDE,
  "longitude": $LONGITUDE,
  "bandwidth_mbps": $BW,
  "max_connections": $MAXC,
  "tags": ["$COUNTRY", "$CITY"]
}
EOF
        cat /opt/nebula/config/node.json
CONFIG
    log "Node configuration created"
}

register_with_backend() {
    if [ -z "${API_URL:-}" ] || [ "$API_URL" = "https://api.nebula-vpn.io" ]; then
        warn "Backend API URL not set, skipping API registration"
        warn "Node config saved at /opt/nebula/config/node.json"
        warn "Register manually later or set API_URL"
        return 0
    fi
    
    info "Registering node with backend API..."
    local config
    config=$($SSH_CMD "cat /opt/nebula/config/node.json")
    
    local resp
    resp=$(curl -s -X POST "$API_URL/v1/admin/servers" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer ${NODE_TOKEN}" \
        -d "$config" 2>&1) || {
        warn "API registration failed: $resp"
        warn "Node is configured locally, register via admin panel later"
        return 0
    }
    log "Node registered with backend API"
}

deploy_containers() {
    info "Deploying monitoring containers..."
    
    $SSH_CMD bash -s <<'COMPOSE'
        set -e
        cat > /opt/nebula/docker-compose.yml <<'YAML'
version: "3.8"

services:
  node-exporter:
    image: prom/node-exporter:latest
    restart: unless-stopped
    network_mode: host
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    command:
      - '--path.procfs=/host/proc'
      - '--path.sysfs=/host/sys'
      - '--path.rootfs=/rootfs'

  wireguard-exporter:
    image: mindflavor/prometheus_wireguard_exporter:latest
    restart: unless-stopped
    network_mode: host
    environment:
      WIREGUARD_EXPORTER_LISTEN: ":9586"

YAML
        cd /opt/nebula && docker compose pull && docker compose up -d
COMPOSE
    log "Monitoring containers deployed"
}

create_health_script() {
    info "Creating health check script..."
    $SSH_CMD bash -s <<'HEALTH'
        set -e
        cat > /opt/nebula/health-check.sh <<'SCRIPT'
#!/usr/bin/env bash
# Simple health check + heartbeat for Nebula VPN node
# Runs every 60s via cron

NODE_CONFIG="/opt/nebula/config/node.json"
WG_DEVICE="$(wg show | head -1 | awk '{print $2}' 2>/dev/null || echo 'none')"
LOAD="$(awk '{print int($1*100/$2)}' /proc/loadavg 2>/dev/null || echo 0)"
MEM_TOTAL="$(awk '/MemTotal/{print $2}' /proc/meminfo 2>/dev/null || echo 1)"
MEM_AVAIL="$(awk '/MemAvailable/{print $2}' /proc/meminfo 2>/dev/null || echo 1)"
MEM_PCT=$(( (MEM_TOTAL - MEM_AVAIL) * 100 / MEM_TOTAL ))
CONNS="$(wg show 2>/dev/null | grep -c "peer:" || echo 0)"

echo "{\"node\": \"$(hostname)\", \"connections\": $CONNS, \"load_pct\": $LOAD, \"mem_pct\": $MEM_PCT, \"ts\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}"
SCRIPT
        chmod +x /opt/nebula/health-check.sh
        
        # Add to cron
        (crontab -l 2>/dev/null; echo "*/5 * * * * /opt/nebula/health-check.sh >> /var/log/nebula/health.log 2>&1") | crontab -
HEALTH
    log "Health check script created (every 5 minutes)"
}

verify_setup() {
    info "Verifying setup..."
    local results
    
    results=$($SSH_CMD bash -s <<'VERIFY'
        echo "--- Docker ---"
        docker info --format '{{.ServerVersion}}' 2>/dev/null || echo "NOT INSTALLED"
        echo "--- WireGuard ---"
        wg --version 2>/dev/null || echo "NOT INSTALLED"
        echo "--- Containers ---"
        docker ps --format 'table {{.Names}}\t{{.Status}}' 2>/dev/null
        echo "--- Firewall ---"
        ufw status verbose 2>/dev/null | head -5
        echo "--- Node Config ---"
        cat /opt/nebula/config/node.json 2>/dev/null || echo "NOT FOUND"
        echo "--- Ports ---"
        ss -tulpn | grep -E ':(51820|9100|9586)' 2>/dev/null || echo "NOT LISTENING"
VERIFY
    )
    echo "$results"
}

# ── Main ──
echo -e "${CYAN}══════════════════════════════════════════${NC}"
echo -e "${CYAN}  Nebula VPN Node Provisioner${NC}"
echo -e "${CYAN}  Server: $SERVER_IP${NC}"
echo -e "${CYAN}══════════════════════════════════════════${NC}"

[ -z "$SERVER_IP" ] && err "Usage: $0 <server-ip> [ssh-port] [ssh-user]"

check_prereqs

echo ""
info "Phase 1/8: Setting up SSH access..."
setup_ssh_key

echo ""
info "Phase 2/8: Checking OS..."
check_os

echo ""
info "Phase 3/8: Installing Docker..."
install_docker

echo ""
info "Phase 4/8: Installing WireGuard..."
install_wireguard

echo ""
info "Phase 5/8: Generating WireGuard keys..."
generate_wireguard_keys

echo ""
info "Phase 6/8: Configuring firewall..."
configure_firewall

echo ""
info "Phase 7/8: Setting up node config..."
setup_node_config

echo ""
info "Phase 8/8: Deploying services..."
deploy_containers
create_health_script
register_with_backend

echo ""
echo -e "${GREEN}══════════════════════════════════════════${NC}"
echo -e "${GREEN}  Server $SERVER_IP provisioned!${NC}"
echo -e "${GREEN}══════════════════════════════════════════${NC}"
echo ""
verify_setup
echo ""
info "Next steps:"
info "1. Add this server to your backend database:"
info "   Server public key: $SERVER_PUBKEY"
info "2. Register in admin panel or via API"
info "3. Clients will see this server in their server list"

# Save server info for later use
mkdir -p /tmp/nebula-provision
cat > "/tmp/nebula-provision/${SERVER_IP}.info" <<EOF
IP=$SERVER_IP
SSH_PORT=$SSH_PORT
COUNTRY=$COUNTRY
CITY=$CITY
PUBKEY=$SERVER_PUBKEY
EOF
