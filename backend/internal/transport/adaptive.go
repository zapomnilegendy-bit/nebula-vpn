package transport

import (
	"context"
	"log/slog"
	"math"
	"sync"
	"time"
)

type ScoreWeights struct {
	Latency    float64
	PacketLoss float64
	Throughput float64
	Jitter     float64
}

type AdaptiveConfig struct {
	Enabled           bool
	CheckInterval     time.Duration
	ScoreWeights      ScoreWeights
	MinStableDuration time.Duration
	FallbackCooldown  time.Duration
}

func DefaultAdaptiveConfig() AdaptiveConfig {
	return AdaptiveConfig{
		Enabled:           true,
		CheckInterval:     30 * time.Second,
		MinStableDuration: 120 * time.Second,
		FallbackCooldown:  60 * time.Second,
		ScoreWeights: ScoreWeights{
			Latency:    0.35,
			PacketLoss: 0.30,
			Throughput: 0.20,
			Jitter:     0.15,
		},
	}
}

type AdaptiveRouter struct {
	transports map[Protocol]Transport
	metrics    map[Protocol]*TransportMetrics
	mu         sync.RWMutex
	current    Protocol
	config     AdaptiveConfig
	lastSwitch time.Time
}

func NewAdaptiveRouter(config AdaptiveConfig) *AdaptiveRouter {
	return &AdaptiveRouter{
		transports: make(map[Protocol]Transport),
		metrics:    make(map[Protocol]*TransportMetrics),
		config:     config,
		current:    WireGuard,
		lastSwitch: time.Now(),
	}
}

func (ar *AdaptiveRouter) RegisterTransport(t Transport) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	ar.transports[t.Name()] = t
}

func (ar *AdaptiveRouter) Start(ctx context.Context) {
	if !ar.config.Enabled {
		return
	}

	ticker := time.NewTicker(ar.config.CheckInterval)
	defer ticker.Stop()

	// Initial evaluation
	ar.evaluateAndSwitch()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ar.evaluateAndSwitch()
		}
	}
}

func (ar *AdaptiveRouter) evaluateAndSwitch() {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	if len(ar.metrics) == 0 {
		return
	}

	// Prevent rapid switching
	if time.Since(ar.lastSwitch) < ar.config.FallbackCooldown {
		return
	}

	if metrics, ok := ar.metrics[ar.current]; ok {
		if time.Since(metrics.Timestamp) < ar.config.MinStableDuration {
			return
		}
	}

	scores := make(map[Protocol]float64)
	for proto, metrics := range ar.metrics {
		scores[proto] = ar.calculateScore(metrics)
	}

	best := ar.current
	bestScore := scores[best]
	for proto, score := range scores {
		if score > bestScore {
			best = proto
			bestScore = score
		}
	}

	if best != ar.current && bestScore > 0 {
		slog.Info("adaptive transport switch",
			"from", ar.current,
			"to", best,
			"score", bestScore,
		)
		ar.current = best
		ar.lastSwitch = time.Now()
	}
}

func (ar *AdaptiveRouter) calculateScore(m *TransportMetrics) float64 {
	if m == nil {
		return 0
	}

	latencyScore := 1.0 / (1.0 + m.RTT.Seconds())
	lossScore := 1.0 - m.PacketLoss
	throughputScore := math.Log2(1.0 + float64(m.ThroughputBps)/1e6)
	jitterScore := 1.0 / (1.0 + m.Jitter.Seconds())

	return ar.config.ScoreWeights.Latency*latencyScore +
		ar.config.ScoreWeights.PacketLoss*lossScore +
		ar.config.ScoreWeights.Throughput*throughputScore +
		ar.config.ScoreWeights.Jitter*jitterScore
}

func (ar *AdaptiveRouter) CurrentTransport() Transport {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	return ar.transports[ar.current]
}

func (ar *AdaptiveRouter) CurrentProtocol() Protocol {
	ar.mu.RLock()
	defer ar.mu.RUnlock()
	return ar.current
}

func (ar *AdaptiveRouter) ReportMetrics(metrics *TransportMetrics) {
	ar.mu.Lock()
	defer ar.mu.Unlock()
	ar.metrics[metrics.Protocol] = metrics
}

func (ar *AdaptiveRouter) AvailableProtocols() []Protocol {
	ar.mu.RLock()
	defer ar.mu.RUnlock()

	protocols := make([]Protocol, 0, len(ar.transports))
	for proto := range ar.transports {
		protocols = append(protocols, proto)
	}
	return protocols
}

func (ar *AdaptiveRouter) ForceSwitch(protocol Protocol) {
	ar.mu.Lock()
	defer ar.mu.Unlock()

	if _, exists := ar.transports[protocol]; exists {
		ar.current = protocol
		ar.lastSwitch = time.Now()
		slog.Info("forced transport switch", "to", protocol)
	}
}
