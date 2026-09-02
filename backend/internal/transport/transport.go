package transport

import (
	"context"
	"net"
	"time"
)

type Protocol string

const (
	WireGuard  Protocol = "wireguard"
	VLESSPlus  Protocol = "vless_reality"
	Hysteria2  Protocol = "hysteria2"
	TUIC       Protocol = "tuic"
)

type ConnConfig struct {
	Protocol   Protocol
	ServerAddr string
	ServerPort int
	PublicKey  string
	PrivateKey string
	MTU        int
	KeepAlive  time.Duration
	Extra      map[string]any
}

type Transport interface {
	Name() Protocol
	Connect(ctx context.Context, config *ConnConfig) (net.Conn, error)
	Disconnect(conn net.Conn) error
	IsHealthy(addr string) bool
	Latency(addr string) (time.Duration, error)
	Metrics() TransportMetrics
}

type TransportMetrics struct {
	Protocol      Protocol
	RTT           time.Duration
	PacketLoss    float64
	ThroughputBps int64
	Jitter        time.Duration
	Timestamp     time.Time
}

type ConnState int

const (
	StateDisconnected ConnState = iota
	StateConnecting
	StateConnected
	StateReconnecting
	StateError
)

func (s ConnState) String() string {
	switch s {
	case StateDisconnected:
		return "disconnected"
	case StateConnecting:
		return "connecting"
	case StateConnected:
		return "connected"
	case StateReconnecting:
		return "reconnecting"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}

type ConnectionInfo struct {
	State        ConnState
	Protocol     Protocol
	ServerAddr   string
	ConnectedAt  time.Time
	BytesRX      int64
	BytesTX      int64
	Latency      time.Duration
}

type ConnectionManager interface {
	Connect(ctx context.Context, config *ConnConfig) error
	Disconnect(ctx context.Context) error
	GetState() ConnectionInfo
	SwitchProtocol(ctx context.Context, protocol Protocol) error
	SetSplitTunnel(ctx context.Context, enabled bool, routes []string) error
	SetKillSwitch(ctx context.Context, enabled bool) error
}
