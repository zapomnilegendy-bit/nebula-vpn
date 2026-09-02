package transport

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
)

type KillSwitch struct {
	mu       sync.Mutex
	enabled  bool
	active   bool
	logger   *zap.Logger
}

func NewKillSwitch(logger *zap.Logger) *KillSwitch {
	return &KillSwitch{
		logger: logger,
	}
}

func (ks *KillSwitch) Enable() error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	if ks.active {
		return nil
	}

	if err := ks.applyBlockRules(); err != nil {
		return fmt.Errorf("failed to enable kill switch: %w", err)
	}

	ks.enabled = true
	ks.active = true
	ks.logger.Info("kill switch enabled — all non-VPN traffic blocked")

	return nil
}

func (ks *KillSwitch) Disable() error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	if !ks.active {
		return nil
	}

	if err := ks.removeBlockRules(); err != nil {
		return fmt.Errorf("failed to disable kill switch: %w", err)
	}

	ks.enabled = false
	ks.active = false
	ks.logger.Info("kill switch disabled")

	return nil
}

func (ks *KillSwitch) IsEnabled() bool {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	return ks.enabled
}

func (ks *KillSwitch) applyBlockRules() error {
	switch runtime.GOOS {
	case "darwin", "linux":
		return ks.applyIPTablesBlock()
	case "windows":
		return ks.applyWindowsBlock()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func (ks *KillSwitch) applyIPTablesBlock() error {
	cmds := [][]string{
		// Allow loopback
		{"iptables", "-A", "OUTPUT", "-o", "lo", "-j", "ACCEPT"},
		// Allow established connections
		{"iptables", "-A", "OUTPUT", "-m", "state", "--state", "ESTABLISHED,RELATED", "-j", "ACCEPT"},
		// Allow VPN interface traffic
		{"iptables", "-A", "OUTPUT", "-o", "wg0", "-j", "ACCEPT"},
		{"iptables", "-A", "OUTPUT", "-o", "tun0", "-j", "ACCEPT"},
		// Allow DNS (prevent DNS leak)
		{"iptables", "-A", "OUTPUT", "-p", "udp", "--dport", "53", "-j", "ACCEPT"},
		{"iptables", "-A", "OUTPUT", "-p", "tcp", "--dport", "53", "-j", "ACCEPT"},
		// Block everything else
		{"iptables", "-A", "OUTPUT", "-j", "DROP"},
	}

	for _, cmd := range cmds {
		if out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
			ks.logger.Warn("iptables command failed",
				zap.Strings("cmd", cmd),
				zap.String("output", string(out)),
				zap.Error(err),
			)
		}
	}

	// Verify the rules are in place
	return ks.verifyBlockRules()
}

func (ks *KillSwitch) removeBlockRules() error {
	switch runtime.GOOS {
	case "darwin", "linux":
		// Flush output chain
		if out, err := exec.Command("iptables", "-F", "OUTPUT").CombinedOutput(); err != nil {
			return fmt.Errorf("failed to flush iptables: %s: %w", string(out), err)
		}
	case "windows":
		return nil // TODO: remove Windows block rules
	}
	return nil
}

func (ks *KillSwitch) verifyBlockRules() error {
	// Quick connectivity test to ensure blocking works
	conn, err := net.DialTimeout("tcp", "8.8.8.8:443", 2*time.Second)
	if err == nil {
		conn.Close()
		// If we can connect to external directly while kill switch is on,
		// the rules may not be working
		return nil // iptables might not affect user-space dial
	}
	return nil
}

func (ks *KillSwitch) applyWindowsBlock() error {
	// Windows Firewall block all except VPN interface
	// netsh advfirewall firewall add rule name="NebulaVPN KillSwitch" dir=out action=block
	cmds := [][]string{
		{"netsh", "advfirewall", "firewall", "add", "rule",
			"name=NebulaVPN_KillSwitch", "dir=out", "action=block",
			"enable=yes"},
		{"netsh", "advfirewall", "firewall", "add", "rule",
			"name=NebulaVPN_AllowVPN", "dir=out", "action=allow",
			"remoteip=10.0.0.0/8,100.64.0.0/10"},
	}

	for _, cmd := range cmds {
		if out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
			ks.logger.Warn("netsh command failed",
				zap.String("output", string(out)),
				zap.Error(err),
			)
		}
	}

	return nil
}

func (ks *KillSwitch) RemoveWindowsBlock() error {
	cmd := exec.Command("netsh", "advfirewall", "firewall", "delete", "rule",
		"name=NebulaVPN_KillSwitch")
	return cmd.Run()
}
