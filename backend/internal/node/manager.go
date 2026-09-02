package node

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NodeRegistration struct {
	Name           string
	Host           string
	IPv4           string
	IPv6           string
	Port           int
	Protocol       string
	PublicKey      string
	Country        string
	City           string
	Latitude       float64
	Longitude      float64
	BandwidthMbps  int
	MaxConnections int
	Tags           []string
}

type RegisteredNode struct {
	ID             string
	Registration   NodeRegistration
	AccessToken    string
	LastHeartbeat  time.Time
	CurrentConns   int
	CPUPercent     float64
	MemPercent     float64
	BwUsedMbps     int64
}

type NodeManager struct {
	mu    sync.RWMutex
	nodes map[string]*RegisteredNode
	logger *zap.Logger
}

func NewNodeManager(logger *zap.Logger) *NodeManager {
	return &NodeManager{
		nodes:  make(map[string]*RegisteredNode),
		logger: logger,
	}
}

func (nm *NodeManager) Register(ctx context.Context, reg *NodeRegistration) (*RegisteredNode, error) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	id := uuid.New().String()
	token, err := generateNodeToken()
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	node := &RegisteredNode{
		ID:           id,
		Registration: *reg,
		AccessToken:  token,
		LastHeartbeat: time.Now(),
	}

	nm.nodes[id] = node
	nm.logger.Info("node registered",
		zap.String("id", id),
		zap.String("host", reg.Host),
		zap.String("country", reg.Country),
	)

	return node, nil
}

func (nm *NodeManager) Heartbeat(nodeID string, conns int, cpu, mem float64, bw int64) {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	if node, ok := nm.nodes[nodeID]; ok {
		node.LastHeartbeat = time.Now()
		node.CurrentConns = conns
		node.CPUPercent = cpu
		node.MemPercent = mem
		node.BwUsedMbps = bw
	}
}

func (nm *NodeManager) GetNode(nodeID string) *RegisteredNode {
	nm.mu.RLock()
	defer nm.mu.RUnlock()
	return nm.nodes[nodeID]
}

func (nm *NodeManager) ListActive() []*RegisteredNode {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	active := make([]*RegisteredNode, 0, len(nm.nodes))
	cutoff := time.Now().Add(-30 * time.Second)

	for _, n := range nm.nodes {
		if n.LastHeartbeat.After(cutoff) {
			active = append(active, n)
		}
	}
	return active
}

func generateNodeToken() (string, error) {
	return uuid.New().String() + "-" + uuid.New().String(), nil
}
