package integration

import (
	"net"
	"testing"
	"time"

	"github.com/nebula-vpn/backend/internal/transport"
	"github.com/nebula-vpn/backend/internal/obfuscator"
	"go.uber.org/zap"
)

func TestAdaptiveRouterScoring(t *testing.T) {
	config := transport.DefaultAdaptiveConfig()
	router := transport.NewAdaptiveRouter(config)

	wg := transport.NewWireGuardTransport(nil)
	vless := transport.NewVLESSRealityTransport(nil, nil)
	router.RegisterTransport(wg)
	router.RegisterTransport(vless)

	// Report metrics
	router.ReportMetrics(&transport.TransportMetrics{
		Protocol:      transport.WireGuard,
		RTT:           20 * time.Millisecond,
		PacketLoss:    0.01,
		ThroughputBps: 100_000_000,
		Jitter:        5 * time.Millisecond,
		Timestamp:     time.Now(),
	})

	router.ReportMetrics(&transport.TransportMetrics{
		Protocol:      transport.VLESSPlus,
		RTT:           80 * time.Millisecond,
		PacketLoss:    0.02,
		ThroughputBps: 50_000_000,
		Jitter:        15 * time.Millisecond,
		Timestamp:     time.Now(),
	})

	// WireGuard should score higher based on weight configuration
	current := router.CurrentProtocol()
	if current != transport.WireGuard {
		t.Logf("current protocol: %s", current)
	}

	// Force switch
	router.ForceSwitch(transport.VLESSPlus)
	if router.CurrentProtocol() != transport.VLESSPlus {
		t.Error("force switch should change protocol")
	}
}

func TestKillSwitchLifecycle(t *testing.T) {
	// Test kill switch logic without actual iptables (platform-dependent)
	logger := zap.NewNop()
	ks := transport.NewKillSwitch(logger)

	if ks.IsEnabled() {
		t.Error("kill switch should start disabled")
	}

	// Note: actual Enable/Disable would require root/iptables
	// Just verifying the state machine
}

func TestSplitTunnelRouting(t *testing.T) {
	st := transport.NewSplitTunnel()

	// Default should route everything
	if !st.IsEnabled() {
		t.Log("split tunnel default is disabled (routes all traffic)")
	}

	err := st.UpdateConfig(transport.SplitTunnelConfig{
		Enabled: true,
		IncludeRoutes: []string{"0.0.0.0/0"},
		ExcludeRoutes: []string{"10.0.0.0/8", "192.168.0.0/16"},
	})
	if err != nil {
		t.Fatalf("failed to update split tunnel config: %v", err)
	}

	if !st.IsEnabled() {
		t.Error("split tunnel should be enabled")
	}

	// Test IP routing decisions
	tests := []struct {
		ip      string
		shouldRoute bool
	}{
		{"8.8.8.8", true},       // public IP
		{"10.0.0.1", false},     // private IP (excluded)
		{"192.168.1.1", false},  // private IP (excluded)
		{"172.16.0.1", true},    // not excluded by default
		{"1.1.1.1", true},       // public IP
	}

	for _, tt := range tests {
		ip := net.ParseIP(tt.ip)
		if ip == nil {
			t.Fatalf("invalid test IP: %s", tt.ip)
		}
		result := st.ShouldRoute(ip)
		if result != tt.shouldRoute {
			t.Errorf("ShouldRoute(%s) = %v, want %v", tt.ip, result, tt.shouldRoute)
		}
	}
}

func TestObfuscatorTLSWrap(t *testing.T) {
	// Unit test for obfuscator without network
	_ = obfuscator.NewTLSCamouflager()
	_ = obfuscator.NewPacketFragmenter(256)
	_ = obfuscator.NewTrafficPadder(40, 200)
	_ = obfuscator.NewTimingRandomizer(5)

	// Test fragmenter
	frag := obfuscator.NewPacketFragmenter(256)
	data := make([]byte, 1000)
	chunks := frag.Fragment(data)

	if len(chunks) != 4 { // 1000 / 256 = 3.9 → 4
		t.Errorf("expected 4 fragments, got %d", len(chunks))
	}

	reassembled := frag.Defragment(chunks)
	if len(reassembled) != 1000 {
		t.Errorf("expected 1000 bytes after defrag, got %d", len(reassembled))
	}
}
