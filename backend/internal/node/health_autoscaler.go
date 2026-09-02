package node

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type HealthChecker struct {
	mu          sync.Mutex
	nodes       map[string]*NodeHealth
	checkInterval time.Duration
	logger      *zap.Logger
	httpClient  *http.Client
}

type NodeHealth struct {
	NodeID     string
	IsHealthy  bool
	LastCheck  time.Time
	Latency    time.Duration
	ErrorCount int
	CPUPercent float64
	MemPercent float64
}

type AutoScaler struct {
	logger         *zap.Logger
	minNodes       int
	maxNodes       int
	cpuThreshold   float64
	memThreshold   float64
	currentNodes   int
	scalingInProgress bool
	mu             sync.Mutex
}

func NewHealthChecker(checkInterval time.Duration, logger *zap.Logger) *HealthChecker {
	return &HealthChecker{
		nodes:         make(map[string]*NodeHealth),
		checkInterval: checkInterval,
		logger:        logger,
		httpClient:   &http.Client{Timeout: 5 * time.Second},
	}
}

func (hc *HealthChecker) Start(ctx context.Context) {
	ticker := time.NewTicker(hc.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			hc.checkAll()
		}
	}
}

func (hc *HealthChecker) RegisterNode(nodeID, addr string) {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	hc.nodes[nodeID] = &NodeHealth{
		NodeID:    nodeID,
		IsHealthy: true,
		LastCheck: time.Now(),
	}
}

func (hc *HealthChecker) checkAll() {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	for id, health := range hc.nodes {
		start := time.Now()

		// TCP health check
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:8080", id), 3*time.Second)
		if err != nil {
			health.ErrorCount++
			health.IsHealthy = false
			hc.logger.Warn("health check failed", zap.String("node", id), zap.Error(err))
			continue
		}
		conn.Close()

		health.Latency = time.Since(start)
		health.LastCheck = time.Now()
		health.ErrorCount = 0
		health.IsHealthy = true
	}
}

func (hc *HealthChecker) GetHealth(nodeID string) *NodeHealth {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	return hc.nodes[nodeID]
}

func (hc *HealthChecker) HealthyNodes() []string {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	var healthy []string
	for id, h := range hc.nodes {
		if h.IsHealthy {
			healthy = append(healthy, id)
		}
	}
	return healthy
}

func NewAutoScaler(logger *zap.Logger, minNodes, maxNodes int, cpuThreshold, memThreshold float64) *AutoScaler {
	return &AutoScaler{
		logger:       logger,
		minNodes:     minNodes,
		maxNodes:     maxNodes,
		cpuThreshold: cpuThreshold,
		memThreshold: memThreshold,
		currentNodes: minNodes,
	}
}

func (as *AutoScaler) Evaluate(nodes []*NodeHealth) {
	as.mu.Lock()
	defer as.mu.Unlock()

	if as.scalingInProgress {
		return
	}

	var avgCPU, avgMem float64
	for _, n := range nodes {
		avgCPU += n.CPUPercent
		avgMem += n.MemPercent
	}
	if len(nodes) > 0 {
		avgCPU /= float64(len(nodes))
		avgMem /= float64(len(nodes))
	}

	// Scale up
	if avgCPU > as.cpuThreshold || avgMem > as.memThreshold {
		if as.currentNodes < as.maxNodes {
			as.currentNodes++
			as.logger.Info("scaling up",
				zap.Int("nodes", as.currentNodes),
				zap.Float64("avg_cpu", avgCPU),
				zap.Float64("avg_mem", avgMem),
			)
		}
	}

	// Scale down
	if avgCPU < as.cpuThreshold*0.5 && avgMem < as.memThreshold*0.5 {
		if as.currentNodes > as.minNodes {
			as.currentNodes--
			as.logger.Info("scaling down",
				zap.Int("nodes", as.currentNodes),
				zap.Float64("avg_cpu", avgCPU),
				zap.Float64("avg_mem", avgMem),
			)
		}
	}
}

func (as *AutoScaler) DesiredNodes() int {
	as.mu.Lock()
	defer as.mu.Unlock()
	return as.currentNodes
}
