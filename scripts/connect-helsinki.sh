#!/usr/bin/env bash
# Nebula VPN — Connect to Helsinki test server
# Usage: sudo ./connect-helsinki.sh

set -e

CONFIG="/tmp/nebula-helsinki.conf"

# Create config
cat > "$CONFIG" <<EOF
[Interface]
PrivateKey = $(cat /tmp/wg-client-priv.key | tr -d '\n')
Address = 10.0.0.2/24
DNS = 1.1.1.1, 8.8.8.8

[Peer]
PublicKey = x4jFp1LCFM1PbixnflU0I6GJJjne5IaN5uyZLxOmdHE=
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = 83.147.216.161:51820
PersistentKeepalive = 25
EOF

echo "Starting WireGuard connection to Helsinki..."
wg-quick up "$CONFIG"

echo ""
echo "=== Connected! ==="
echo "Test: ping 10.0.0.1"
echo "Test: curl ifconfig.me"
echo ""
echo "To disconnect: wg-quick down $CONFIG"
