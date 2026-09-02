package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

type VLESSRealityTransport struct {
	logger       *zap.Logger
	serverNames  []string
	tlsConfig    *tls.Config
}

func NewVLESSRealityTransport(logger *zap.Logger, serverNames []string) *VLESSRealityTransport {
	if len(serverNames) == 0 {
		serverNames = []string{"www.microsoft.com", "www.apple.com", "www.cloudflare.com"}
	}

	return &VLESSRealityTransport{
		logger:      logger,
		serverNames: serverNames,
		tlsConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
			MaxVersion: tls.VersionTLS13,
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
			},
			InsecureSkipVerify: true,
		},
	}
}

func (v *VLESSRealityTransport) Name() Protocol {
	return VLESSPlus
}

func (v *VLESSRealityTransport) Connect(ctx context.Context, config *ConnConfig) error {
	addr := fmt.Sprintf("%s:%d", config.ServerAddr, config.ServerPort)

	dialer := &net.Dialer{Timeout: 10 * time.Second}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("vless+tcp dial failed: %w", err)
	}

	// Wrap in TLS with reality-like server name
	v.tlsConfig.ServerName = v.serverNames[randInt(len(v.serverNames))]
	tlsConn := tls.Client(conn, v.tlsConfig)

	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		return fmt.Errorf("vless+tls handshake failed: %w", err)
	}

	v.logger.Info("VLESS+Reality connection established",
		zap.String("endpoint", addr),
		zap.String("sni", v.tlsConfig.ServerName),
	)

	return nil
}

func (v *VLESSRealityTransport) Disconnect(conn net.Conn) error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func (v *VLESSRealityTransport) IsHealthy(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (v *VLESSRealityTransport) Latency(addr string) (time.Duration, error) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return 0, fmt.Errorf("tcp dial failed: %w", err)
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, _ = conn.Write([]byte{0}) // VLESS handshake start

	duration := time.Since(start)
	return duration, nil
}

func (v *VLESSRealityTransport) Metrics() TransportMetrics {
	return TransportMetrics{
		Protocol:  VLESSPlus,
		Timestamp: time.Now(),
	}
}
