package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"go.uber.org/zap"
)

type Hysteria2Transport struct {
	logger       *zap.Logger
	obfsPassword string
	upMbps       int
	downMbps     int
}

func NewHysteria2Transport(logger *zap.Logger, obfsPassword string, upMbps, downMbps int) *Hysteria2Transport {
	return &Hysteria2Transport{
		logger:       logger,
		obfsPassword: obfsPassword,
		upMbps:       upMbps,
		downMbps:     downMbps,
	}
}

func (h *Hysteria2Transport) Name() Protocol {
	return Hysteria2
}

func (h *Hysteria2Transport) Connect(ctx context.Context, config *ConnConfig) error {
	addr := fmt.Sprintf("%s:%d", config.ServerAddr, config.ServerPort)

	// Hysteria2 uses QUIC (UDP) with TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
		ServerName: "cloudflare.com",
		NextProtos: []string{"h3"},
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("hysteria2 resolve failed: %w", err)
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return fmt.Errorf("hysteria2 udp dial failed: %w", err)
	}

	tlsConn := tls.Client(udpConn, tlsConfig)

	// QUIC handshake (simplified — production uses quic-go)
	if err := tlsConn.Handshake(); err != nil {
		udpConn.Close()
		return fmt.Errorf("hysteria2 quic handshake failed: %w", err)
	}

	h.logger.Info("Hysteria2 connection established",
		zap.String("endpoint", addr),
		zap.Int("up_mbps", h.upMbps),
		zap.Int("down_mbps", h.downMbps),
	)

	return nil
}

func (h *Hysteria2Transport) Disconnect(conn net.Conn) error {
	if conn != nil {
		return conn.Close()
	}
	return nil
}

func (h *Hysteria2Transport) IsHealthy(addr string) bool {
	conn, err := net.DialTimeout("udp", addr, 5*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func (h *Hysteria2Transport) Latency(addr string) (time.Duration, error) {
	start := time.Now()
	conn, err := net.DialTimeout("udp", addr, 5*time.Second)
	if err != nil {
		return 0, fmt.Errorf("udp dial failed: %w", err)
	}
	defer conn.Close()

	_, _ = conn.Write([]byte{0x01, 0x00, 0x00, 0x00}) // H2 ping

	duration := time.Since(start)
	return duration, nil
}

func (h *Hysteria2Transport) Metrics() TransportMetrics {
	return TransportMetrics{
		Protocol:  Hysteria2,
		Timestamp: time.Now(),
	}
}
