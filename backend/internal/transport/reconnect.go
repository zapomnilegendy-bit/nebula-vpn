package transport

import (
	"context"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ReconnectManager struct {
	logger        *zap.Logger
	config        ReconnectConfig
	mu            sync.Mutex
	state         ConnState
	retryCount    int
	activeConn    context.CancelFunc
	reconnectCh   chan struct{}
	doneCh        chan struct{}
	endpoints     []FailoverEndpoint
	currentEp     int
}

type ReconnectConfig struct {
	MaxRetries        int
	BaseBackoff       time.Duration
	MaxBackoff        time.Duration
	JitterPercent     float64
	HealthCheckInterval time.Duration
}

type FailoverEndpoint struct {
	Host    string
	Port    int
	Country string
	Weight  int
}

func DefaultReconnectConfig() ReconnectConfig {
	return ReconnectConfig{
		MaxRetries:         10,
		BaseBackoff:        time.Second,
		MaxBackoff:         60 * time.Second,
		JitterPercent:      0.2,
		HealthCheckInterval: 15 * time.Second,
	}
}

func NewReconnectManager(logger *zap.Logger, config ReconnectConfig) *ReconnectManager {
	return &ReconnectManager{
		logger:      logger,
		config:      config,
		state:       StateDisconnected,
		reconnectCh: make(chan struct{}, 1),
		doneCh:      make(chan struct{}),
	}
}

func (rm *ReconnectManager) SetEndpoints(endpoints []FailoverEndpoint) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.endpoints = endpoints
	rm.currentEp = 0
}

func (rm *ReconnectManager) Start(ctx context.Context, connectFn func(ctx context.Context, ep FailoverEndpoint) error) {
	go rm.loop(ctx, connectFn)
}

func (rm *ReconnectManager) loop(ctx context.Context, connectFn func(ctx context.Context, ep FailoverEndpoint) error) {
	healthTicker := time.NewTicker(rm.config.HealthCheckInterval)
	defer healthTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-rm.reconnectCh:
			rm.mu.Lock()
			rm.state = StateReconnecting
			rm.mu.Unlock()

			rm.tryConnectWithBackoff(ctx, connectFn)

		case <-healthTicker.C:
			rm.mu.Lock()
			shouldCheck := rm.state == StateConnected
			rm.mu.Unlock()

			if shouldCheck {
				rm.checkAndRecover(ctx, connectFn)
			}
		}
	}
}

func (rm *ReconnectManager) tryConnectWithBackoff(ctx context.Context, connectFn func(ctx context.Context, ep FailoverEndpoint) error) {
	for i := 0; i < rm.config.MaxRetries; i++ {
		select {
		case <-ctx.Done():
			return
		default:
		}

		ep := rm.selectEndpoint()

		backoff := rm.calculateBackoff(i)
		zap.L().Info("reconnect attempt",
			zap.Int("attempt", i+1),
			zap.Int("max", rm.config.MaxRetries),
			zap.String("endpoint", fmt.Sprintf("%s:%d", ep.Host, ep.Port)),
			zap.Duration("backoff", backoff),
		)

		connectCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		err := connectFn(connectCtx, ep)
		cancel()

		if err == nil {
			rm.mu.Lock()
			rm.state = StateConnected
			rm.retryCount = 0
			rm.mu.Unlock()

			zap.L().Info("reconnect successful",
				zap.String("endpoint", fmt.Sprintf("%s:%d", ep.Host, ep.Port)),
				zap.Int("attempts", i+1),
			)
			return
		}

		zap.L().Warn("reconnect failed",
			zap.Int("attempt", i+1),
			zap.Error(err),
		)

		// Failover to next endpoint
		rm.mu.Lock()
		rm.currentEp = (rm.currentEp + 1) % max(1, len(rm.endpoints))
		rm.mu.Unlock()

		// Wait with backoff (but respect context cancellation)
		select {
		case <-ctx.Done():
			return
		case <-time.After(backoff):
		}
	}

	rm.mu.Lock()
	rm.state = StateError
	rm.mu.Unlock()

	zap.L().Error("all reconnect attempts exhausted")
}

func (rm *ReconnectManager) checkAndRecover(ctx context.Context, connectFn func(ctx context.Context, ep FailoverEndpoint) error) {
	ep := rm.selectEndpoint()

	if !isReachable(ep.Host, ep.Port) {
		zap.L().Warn("endpoint unreachable, triggering reconnect",
			zap.String("endpoint", fmt.Sprintf("%s:%d", ep.Host, ep.Port)),
		)
		rm.TriggerReconnect()
	}
}

func (rm *ReconnectManager) TriggerReconnect() {
	select {
	case rm.reconnectCh <- struct{}{}:
	default:
	}
}

func (rm *ReconnectManager) selectEndpoint() FailoverEndpoint {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if len(rm.endpoints) == 0 {
		return FailoverEndpoint{Host: "127.0.0.1", Port: 51820}
	}

	idx := rm.currentEp % len(rm.endpoints)
	return rm.endpoints[idx]
}

func (rm *ReconnectManager) calculateBackoff(attempt int) time.Duration {
	backoff := float64(rm.config.BaseBackoff) * math.Pow(2, float64(attempt))
	if backoff > float64(rm.config.MaxBackoff) {
		backoff = float64(rm.config.MaxBackoff)
	}

	// Add jitter (±jitterPercent)
	jitter := backoff * rm.config.JitterPercent
	backoff += jitter * (randFloat64()*2 - 1)

	return time.Duration(backoff)
}

func isReachable(host string, port int) bool {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func randFloat64() float64 {
	return float64(time.Now().UnixNano()%1000) / 1000.0
}
