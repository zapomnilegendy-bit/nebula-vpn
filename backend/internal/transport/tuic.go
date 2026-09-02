package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

type TUICTransport struct {
	logger    *zap.Logger
	tlsConfig *tls.Config
}

func NewTUICTransport(logger *zap.Logger) *TUICTransport {
	return &TUICTransport{
		logger: logger,
		tlsConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
			MaxVersion: tls.VersionTLS13,
			ServerName: "relay.cloudflare.com",
			NextProtos: []string{"h3", "h2"},
		},
	}
}

func (t *TUICTransport) Name() Protocol {
	return TUIC
}

func (t *TUICTransport) Connect(ctx context.Context, config *ConnConfig) error {
	addr := fmt.Sprintf("%s:%d", config.ServerAddr, config.ServerPort)

	// TUIC uses TCP+TLS+QUIC multiplexing
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	tcpConn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("tuic tcp dial failed: %w", err)
	}

	tlsConn := tls.Client(tcpConn, t.tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		tcpConn.Close()
		return fmt.Errorf("tuic tls handshake failed: %w", err)
	}

	t.logger.Info("TUIC connection established",
		zap.String("endpoint", addr),
		zap.String("alpn", "h3/h2"),
	)

	return nil
}

func (t *TUICTransport) Disconnect(conn net.Conn) error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func (t *TUICTransport) IsHealthy(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (t *TUICTransport) Latency(addr string) (time.Duration, error) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return 0, fmt.Errorf("tcp dial failed: %w", err)
	}
	defer conn.Close()

	duration := time.Since(start)
	return duration, nil
}

func (t *TUICTransport) Metrics() TransportMetrics {
	return TransportMetrics{
		Protocol:  TUIC,
		Timestamp: time.Now(),
	}
}
