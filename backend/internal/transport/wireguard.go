package transport

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"go.uber.org/zap"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type WireGuardTransport struct {
	logger *zap.Logger
	client *wgctrl.Client
}

func NewWireGuardTransport(logger *zap.Logger) *WireGuardTransport {
	client, err := wgctrl.New()
	if err != nil {
		logger.Warn("failed to create wgctrl client, using exec fallback", zap.Error(err))
	}

	return &WireGuardTransport{
		logger: logger,
		client: client,
	}
}

func (w *WireGuardTransport) Name() Protocol {
	return WireGuard
}

func (w *WireGuardTransport) Connect(ctx context.Context, config *ConnConfig) error {
	// Generate WireGuard keypair if not provided
	privateKey := config.PrivateKey
	publicKey := config.PublicKey

	if privateKey == "" {
		key, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return fmt.Errorf("failed to generate private key: %w", err)
		}
		privateKey = key.String()
		publicKey = key.PublicKey().String()
	}

	// Build config file
	configContent := fmt.Sprintf(`[Interface]
PrivateKey = %s
Address = %s/32
DNS = 1.1.1.1, 1.0.0.1
MTU = %d

[Peer]
PublicKey = %s
Endpoint = %s:%d
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = %d
`,
		privateKey,
		config.Extra["address"],
		config.MTU,
		publicKey,
		config.ServerAddr,
		config.ServerPort,
		int(config.KeepAlive.Seconds()),
	)

	// Write config to temp file
	tmpFile, err := os.CreateTemp("", "wg-*.conf")
	if err != nil {
		return fmt.Errorf("failed to create temp config: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write config: %w", err)
	}
	tmpFile.Close()

	// Bring up interface using wg-quick
	cmd := exec.CommandContext(ctx, "wg-quick", "up", tmpFile.Name())
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to bring up wireguard: %s: %w", string(output), err)
	}

	w.logger.Info("wireguard connection established",
		zap.String("endpoint", fmt.Sprintf("%s:%d", config.ServerAddr, config.ServerPort)),
		zap.Int("mtu", config.MTU),
	)

	return nil
}

func (w *WireGuardTransport) Disconnect(conn net.Conn) error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func (w *WireGuardTransport) IsHealthy(addr string) bool {
	conn, err := net.DialTimeout("udp", addr, 5*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (w *WireGuardTransport) Latency(addr string) (time.Duration, error) {
	start := time.Now()
	conn, err := net.DialTimeout("udp", addr, 5*time.Second)
	if err != nil {
		return 0, fmt.Errorf("connect failed: %w", err)
	}
	defer conn.Close()

	// Simple handshake latency
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_ = conn.SetWriteDeadline(time.Now().Add(3 * time.Second))

	duration := time.Since(start)
	return duration, nil
}

func (w *WireGuardTransport) Metrics() TransportMetrics {
	// Query wgctrl for device stats
	if w.client != nil {
		devices, err := w.client.Devices()
		if err == nil && len(devices) > 0 {
			dev := devices[0]
			var totalRX, totalTX int64
			for _, peer := range dev.Peers {
				totalRX += int64(peer.ReceiveBytes)
				totalTX += int64(peer.TransmitBytes)
			}

			return TransportMetrics{
				Protocol:      WireGuard,
				ThroughputBps: totalRX + totalTX,
				Timestamp:     time.Now(),
			}
		}
	}

	return TransportMetrics{
		Protocol:  WireGuard,
		Timestamp: time.Now(),
	}
}

func (w *WireGuardTransport) GenerateKeypair() (*wgtypes.Keypair, error) {
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate keypair: %w", err)
	}

	return &wgtypes.Keypair{
		PrivateKey: privateKey,
		PublicKey:  privateKey.PublicKey(),
	}, nil
}

func (w *WireGuardTransport) ParseConfig(configStr string) (*ConnConfig, error) {
	// Parse standard WireGuard config format
	cfg := &ConnConfig{
		Protocol: WireGuard,
		Extra:    make(map[string]any),
	}

	lines := strings.Split(configStr, "\n")
	var currentSection string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.ToLower(line[1 : len(line)-1])
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch currentSection {
		case "interface":
			switch key {
			case "PrivateKey":
				cfg.PrivateKey = value
			case "Address":
				cfg.Extra["address"] = value
			case "MTU":
				fmt.Sscanf(value, "%d", &cfg.MTU)
			case "DNS":
				cfg.Extra["dns"] = value
			}
		case "peer":
			switch key {
			case "PublicKey":
				cfg.PublicKey = value
			case "Endpoint":
				host, portStr, _ := net.SplitHostPort(value)
				cfg.ServerAddr = host
				fmt.Sscanf(portStr, "%d", &cfg.ServerPort)
			case "PersistentKeepalive":
				var keepalive int
				fmt.Sscanf(value, "%d", &keepalive)
				cfg.KeepAlive = time.Duration(keepalive) * time.Second
			}
		}
	}

	return cfg, nil
}
