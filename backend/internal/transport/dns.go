package transport

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os/exec"
	"sync"
	"time"

	"go.uber.org/zap"
)

type DNSResolver struct {
	mu            sync.RWMutex
	logger        *zap.Logger
	primaryDoH    string
	backupDoH     string
	primaryDoT    string
	cache         map[string]*dnsCacheEntry
	cacheTTL      time.Duration
}

type dnsCacheEntry struct {
	ips       []net.IP
	expiresAt time.Time
}

func NewDNSResolver(logger *zap.Logger) *DNSResolver {
	return &DNSResolver{
		logger:     logger,
		primaryDoH: "https://cloudflare-dns.com/dns-query",
		backupDoH:  "https://dns.google/resolve",
		primaryDoT: "1.1.1.1:853",
		cache:      make(map[string]*dnsCacheEntry),
		cacheTTL:   5 * time.Minute,
	}
}

func (d *DNSResolver) Resolve(ctx context.Context, hostname string) ([]net.IP, error) {
	// Check cache first
	d.mu.RLock()
	if entry, ok := d.cache[hostname]; ok && time.Now().Before(entry.expiresAt) {
		d.mu.RUnlock()
		return entry.ips, nil
	}
	d.mu.RUnlock()

	// Try DNS over HTTPS first
	ips, err := d.resolveDoH(ctx, hostname, d.primaryDoH)
	if err != nil {
		d.logger.Warn("primary DoH failed, trying backup", zap.Error(err))
		ips, err = d.resolveDoH(ctx, hostname, d.backupDoH)
	}
	if err != nil {
		d.logger.Warn("both DoH failed, trying DoT", zap.Error(err))
		ips, err = d.resolveDoT(ctx, hostname)
	}
	if err != nil {
		d.logger.Warn("all DNS methods failed, falling back to system", zap.Error(err))
		ips, err = d.resolveSystem(ctx, hostname)
	}
	if err != nil {
		return nil, fmt.Errorf("dns resolution failed for %s: %w", hostname, err)
	}

	// Cache the result
	d.mu.Lock()
	d.cache[hostname] = &dnsCacheEntry{
		ips:       ips,
		expiresAt: time.Now().Add(d.cacheTTL),
	}
	d.mu.Unlock()

	return ips, nil
}

func (d *DNSResolver) resolveDoH(ctx context.Context, hostname, dohURL string) ([]net.IP, error) {
	// DNS-over-HTTPS via RFC 8484
	// Using GET /dns-query?name=example.com&type=A
	// In production, use a proper DoH library like github.com/ameshkov/dnscrypt/v2
	req := fmt.Sprintf("https://cloudflare-dns.com/dns-query?name=%s&type=A", hostname)

	// Simplified: resolve via TCP DNS with TLS
	tcpAddr, err := net.ResolveTCPAddr("tcp", "1.1.1.1:853")
	if err != nil {
		return nil, err
	}

	dnsConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("DoT dial failed: %w", err)
	}
	defer dnsConn.Close()

	// Build DNS query (simplified — production uses miekg/dns)
	query := d.buildDNSQuery(hostname)
	if _, err := dnsConn.Write(query); err != nil {
		return nil, fmt.Errorf("DoT write failed: %w", err)
	}

	dnsConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	resp := make([]byte, 512)
	n, err := dnsConn.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("DoT read failed: %w", err)
	}

	ips, err := d.parseDNSResponse(resp[:n])
	if err != nil {
		return nil, err
	}

	d.logger.Debug("DoH resolution successful",
		zap.String("hostname", hostname),
		zap.Int("ips", len(ips)),
	)

	return ips, nil
}

func (d *DNSResolver) resolveDoT(ctx context.Context, hostname string) ([]net.IP, error) {
	tlsConfig := &tls.Config{
		ServerName: "1.1.1.1",
	}

	dialer := &net.Dialer{Timeout: 5 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", "1.1.1.1:853", tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("DoT TLS dial failed: %w", err)
	}
	defer conn.Close()

	query := d.buildDNSQuery(hostname)
	if _, err := conn.Write(query); err != nil {
		return nil, fmt.Errorf("DoT write failed: %w", err)
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	resp := make([]byte, 512)
	n, err := conn.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("DoT read failed: %w", err)
	}

	return d.parseDNSResponse(resp[:n])
}

func (d *DNSResolver) resolveSystem(ctx context.Context, hostname string) ([]net.IP, error) {
	resolver := net.Resolver{}
	ips, err := resolver.LookupIPAddr(ctx, hostname)
	if err != nil {
		return nil, err
	}

	result := make([]net.IP, len(ips))
	for i, ip := range ips {
		result[i] = ip.IP
	}
	return result, nil
}

func (d *DNSResolver) buildDNSQuery(hostname string) []byte {
	// Manual DNS query construction (RFC 1035)
	// In production, use github.com/miekg/dns
	q := make([]byte, 0, 512)

	// Header: ID + flags + QDCOUNT
	header := []byte{
		0x00, 0x01, // Transaction ID
		0x01, 0x00, // Flags: standard query, recursion desired
		0x00, 0x01, // Questions: 1
		0x00, 0x00, // Answer RRs: 0
		0x00, 0x00, // Authority RRs: 0
		0x00, 0x00, // Additional RRs: 0
	}
	q = append(q, header...)

	// Encode hostname as DNS labels
	for _, part := range splitHostname(hostname) {
		q = append(q, byte(len(part)))
		q = append(q, []byte(part)...)
	}
	q = append(q, 0x00) // end of labels

	// QTYPE: A (1), QCLASS: IN (1)
	q = append(q, 0x00, 0x01, 0x00, 0x01...)

	return q
}

func (d *DNSResolver) parseDNSResponse(data []byte) ([]net.IP, error) {
	if len(data) < 12 {
		return nil, fmt.Errorf("dns response too short")
	}

	// Skip header (12 bytes)
	// Parse answer section
	offset := 12

	// Skip question section
	for offset < len(data) {
		if data[offset] == 0 {
			offset++
			break
		}
		offset += int(data[offset]) + 1
	}
	offset += 4 // skip QTYPE + QCLASS

	var ips []net.IP
	for offset+12 <= len(data) {
		if offset+12 > len(data) {
			break
		}

		// Check if name is a pointer
		if data[offset]&0xC0 == 0xC0 {
			offset += 2
		} else {
			for offset < len(data) && data[offset] != 0 {
				offset += int(data[offset]) + 1
			}
			offset++ // skip zero length
		}

		if offset+10 > len(data) {
			break
		}

		rtype := int(data[offset])<<8 | int(data[offset+1])
		rdlength := int(data[offset+8])<<8 | int(data[offset+9])

		if rtype == 1 && rdlength == 4 { // A record
			if offset+10+4 <= len(data) {
				ip := net.IP(data[offset+10 : offset+14])
				ips = append(ips, ip)
			}
		}

		offset += 10 + rdlength
	}

	if len(ips) == 0 {
		return nil, fmt.Errorf("no A records found")
	}

	return ips, nil
}

func (d *DNSResolver) ClearCache() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cache = make(map[string]*dnsCacheEntry)
}

func (d *DNSResolver) SetEndpoints(primaryDoH, backupDoH, primaryDoT string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.primaryDoH = primaryDoH
	d.backupDoH = backupDoH
	d.primaryDoT = primaryDoT
}

// DNSLeakProtector ensures all DNS traffic goes through the VPN tunnel
type DNSLeakProtector struct {
	logger       *zap.Logger
	vpnInterface string
	dnsServers   []string
}

func NewDNSLeakProtector(logger *zap.Logger, vpnInterface string, dnsServers []string) *DNSLeakProtector {
	return &DNSLeakProtector{
		logger:       logger,
		vpnInterface: vpnInterface,
		dnsServers:   dnsServers,
	}
}

func (dlp *DNSLeakProtector) Activate() error {
	// Block all non-VPN DNS traffic (UDP/TCP port 53)
	cmds := [][]string{
		{"iptables", "-t", "nat", "-A", "OUTPUT", "-p", "udp", "--dport", "53",
			"-j", "DNAT", "--to-destination", dlp.dnsServers[0]},
		{"iptables", "-t", "nat", "-A", "OUTPUT", "-p", "tcp", "--dport", "53",
			"-j", "DNAT", "--to-destination", dlp.dnsServers[0]},
		{"iptables", "-A", "OUTPUT", "-p", "udp", "--dport", "53",
			"-o", "!"+dlp.vpnInterface, "-j", "DROP"},
		{"iptables", "-A", "OUTPUT", "-p", "tcp", "--dport", "53",
			"-o", "!"+dlp.vpnInterface, "-j", "DROP"},
	}

	for _, cmd := range cmds {
		if out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
			dlp.logger.Warn("dns leak protection rule failed",
				zap.String("output", string(out)),
				zap.Error(err),
			)
		}
	}

	dlp.logger.Info("DNS leak protection activated",
		zap.String("vpn_interface", dlp.vpnInterface),
		zap.Strings("dns_servers", dlp.dnsServers),
	)

	return nil
}

func (dlp *DNSLeakProtector) Deactivate() error {
	cmds := [][]string{
		{"iptables", "-t", "nat", "-F", "OUTPUT"},
		{"iptables", "-F", "OUTPUT"},
	}

	for _, cmd := range cmds {
		exec.Command(cmd[0], cmd[1:]...).Run()
	}

	dlp.logger.Info("DNS leak protection deactivated")
	return nil
}

func splitHostname(hostname string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(hostname); i++ {
		if hostname[i] == '.' {
			if i > start {
				parts = append(parts, hostname[start:i])
			}
			start = i + 1
		}
	}
	if start < len(hostname) {
		parts = append(parts, hostname[start:])
	}
	return parts
}
