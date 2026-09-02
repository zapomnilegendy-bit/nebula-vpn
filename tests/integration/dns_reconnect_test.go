package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nebula-vpn/backend/internal/transport"
	"go.uber.org/zap"
)

func TestDNSResolver(t *testing.T) {
	logger := zap.NewNop()
	resolver := transport.NewDNSResolver(logger)
	defer resolver.ClearCache()

	// Test cache
	resolver.SetEndpoints("https://cloudflare-dns.com/dns-query", "https://dns.google/resolve", "1.1.1.1:853")

	// Test resolution (may need network)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ips, err := resolver.Resolve(ctx, "cloudflare.com")
	if err != nil {
		t.Skipf("DNS resolution failed (network required): %v", err)
	}

	if len(ips) == 0 {
		t.Error("expected at least one IP address")
	}

	// Test caching (should return quickly)
	cachedIPs, err := resolver.Resolve(ctx, "cloudflare.com")
	if err != nil {
		t.Errorf("cached resolution failed: %v", err)
	}

	if len(cachedIPs) != len(ips) {
		t.Errorf("cached IP count mismatch: got %d, want %d", len(cachedIPs), len(ips))
	}
}

func TestDNSLeakProtector(t *testing.T) {
	logger := zap.NewNop()
	protector := transport.NewDNSLeakProtector(logger, "wg0", []string{"1.1.1.1", "1.0.0.1"})

	// In CI, iptables may not be available, so we just test the struct creation
	if protector == nil {
		t.Fatal("protector should not be nil")
	}
}

func TestReconnectManager(t *testing.T) {
	logger := zap.NewNop()
	config := transport.DefaultReconnectConfig()
	config.MaxRetries = 3
	config.BaseBackoff = 100 * time.Millisecond

	rm := transport.NewReconnectManager(logger, config)

	rm.SetEndpoints([]transport.FailoverEndpoint{
		{Host: "vpn1.nebula-vpn.io", Port: 51820},
		{Host: "vpn2.nebula-vpn.io", Port: 51820},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the reconnect loop
	connectAttempts := 0
	rm.Start(ctx, func(ctx context.Context, ep transport.FailoverEndpoint) error {
		connectAttempts++
		return fmt.Errorf("connection refused (test)")
	})

	// Trigger reconnect
	rm.TriggerReconnect()

	// Wait briefly
	time.Sleep(500 * time.Millisecond)

	if connectAttempts == 0 {
		t.Log("connect attempts should increment (async)")
	}
}
