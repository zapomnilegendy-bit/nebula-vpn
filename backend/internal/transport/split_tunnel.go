package transport

import (
	"fmt"
	"net"
	"sync"
)

type SplitTunnelConfig struct {
	Enabled        bool
	IncludeRoutes  []string
	ExcludeRoutes  []string
	IncludeDomains []string
	ExcludeDomains []string
	IncludeIPs     []net.IPNet
	ExcludeIPs     []net.IPNet
}

type SplitTunnel struct {
	mu     sync.RWMutex
	config SplitTunnelConfig
}

func NewSplitTunnel() *SplitTunnel {
	return &SplitTunnel{
		config: SplitTunnelConfig{
			Enabled: false,
			IncludeRoutes: []string{"0.0.0.0/0", "::/0"},
			ExcludeRoutes: []string{
				"10.0.0.0/8",
				"172.16.0.0/12",
				"192.168.0.0/16",
				"100.64.0.0/10",
			},
		},
	}
}

func (st *SplitTunnel) UpdateConfig(config SplitTunnelConfig) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	// Parse and validate all routes
	includeIPs, err := parseCIDRs(config.IncludeRoutes)
	if err != nil {
		return fmt.Errorf("invalid include route: %w", err)
	}
	config.IncludeIPs = includeIPs

	excludeIPs, err := parseCIDRs(config.ExcludeRoutes)
	if err != nil {
		return fmt.Errorf("invalid exclude route: %w", err)
	}
	config.ExcludeIPs = excludeIPs

	if config.Enabled && len(config.IncludeRoutes) == 0 {
		return fmt.Errorf("split tunnel enabled but no include routes specified")
	}

	st.config = config
	return nil
}

func (st *SplitTunnel) ShouldRoute(ip net.IP) bool {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if !st.config.Enabled {
		return true
	}

	if st.config.ExcludeIPs != nil {
		for _, cidr := range st.config.ExcludeIPs {
			if cidr.Contains(ip) {
				return false
			}
		}
	}

	if st.config.IncludeIPs != nil {
		for _, cidr := range st.config.IncludeIPs {
			if cidr.Contains(ip) {
				return true
			}
		}
		return false
	}

	return true
}

func (st *SplitTunnel) GetAllowedIPs() []string {
	st.mu.RLock()
	defer st.mu.RUnlock()

	if !st.config.Enabled {
		return []string{"0.0.0.0/0", "::/0"}
	}

	// Return include routes for the WireGuard config
	if len(st.config.IncludeRoutes) > 0 {
		return st.config.IncludeRoutes
	}

	return []string{"0.0.0.0/0", "::/0"}
}

func (st *SplitTunnel) IsEnabled() bool {
	st.mu.RLock()
	defer st.mu.RUnlock()
	return st.config.Enabled
}

func parseCIDRs(routes []string) ([]net.IPNet, error) {
	parsed := make([]net.IPNet, 0, len(routes))
	for _, route := range routes {
		_, cidr, err := net.ParseCIDR(route)
		if err != nil {
			return nil, fmt.Errorf("invalid CIDR %q: %w", route, err)
		}
		parsed = append(parsed, *cidr)
	}
	return parsed, nil
}
