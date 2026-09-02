package obfuscator

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"math/big"
	"net"
	"time"
)

type TLSCamouflager struct {
	ja3Pool []string
}

func NewTLSCamouflager() *TLSCamouflager {
	return &TLSCamouflager{
		ja3Pool: generateJA3Pool(50),
	}
}

func (tc *TLSCamouflager) Wrap(conn net.Conn) (net.Conn, error) {
	ja3 := tc.ja3Pool[randInt(len(tc.ja3Pool))]

	tlsConn := tls.Client(conn, &tls.Config{
		MinVersion:       tls.VersionTLS13,
		MaxVersion:       tls.VersionTLS13,
		CipherSuites:     ja3ToCiphers(ja3),
		InsecureSkipVerify: true,
		ServerName:       randomSNI(),
	})

	if err := tlsConn.Handshake(); err != nil {
		return nil, err
	}

	return &camouflagedConn{
		Conn: tlsConn,
	}, nil
}

type camouflagedConn struct {
	*tls.Conn
}

func (c *camouflagedConn) Write(data []byte) (int, error) {
	padSize := randPadding()
	padded := make([]byte, len(data)+padSize)
	copy(padded, data)
	rand.Read(padded[len(data):])

	time.Sleep(time.Duration(randInt(10)) * time.Millisecond)

	return c.Conn.Write(padded)
}

func (c *camouflagedConn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	if err != nil {
		return n, err
	}
	return n, nil
}

var ja3Fingerprints = []string{
	"771,4865-4866-4867-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0",
	"771,4865-4867-4866-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0",
	"771,4865-4866-4867-49196-49200-49195-49199-159-160-52393-52392-161-162-49171-49172-163-164-156-157-47-53-49187-49189-49161-49163-255,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21-17513,29-23-24,0",
	"771,4865-4866-4867-49196-49200-49195-49199-159-160-52393-52392-162-163-164-161-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21-17513-2570,29-23-24,0",
	"771,4865-4867-4866-49195-49199-49196-49200-52393-52392-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-18-51-45-43-27-21,29-23-24,0",
}

func generateJA3Pool(count int) []string {
	if count > len(ja3Fingerprints) {
		count = len(ja3Fingerprints)
	}
	return ja3Fingerprints[:count]
}

func ja3ToCiphers(ja3 string) []uint16 {
	return []uint16{
		tls.TLS_AES_128_GCM_SHA256,
		tls.TLS_AES_256_GCM_SHA384,
		tls.TLS_CHACHA20_POLY1305_SHA256,
	}
}

var sniList = []string{
	"www.cloudflare.com",
	"www.microsoft.com",
	"www.apple.com",
	"www.github.com",
	"www.amazon.com",
	"www.google.com",
	"www.youtube.com",
}

func randomSNI() string {
	return sniList[randInt(len(sniList))]
}

func randPadding() int {
	return 40 + randInt(161)
}

func randInt(max int) int {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return int(n.Int64())
}

func randomizeBuffer(buf []byte) {
	rand.Read(buf)
}

type PacketFragmenter struct {
	FragmentSize int
}

func NewPacketFragmenter(fragmentSize int) *PacketFragmenter {
	if fragmentSize == 0 {
		fragmentSize = 256
	}
	return &PacketFragmenter{FragmentSize: fragmentSize}
}

func (pf *PacketFragmenter) Fragment(data []byte) [][]byte {
	if len(data) <= pf.FragmentSize {
		return [][]byte{data}
	}

	var fragments [][]byte
	for i := 0; i < len(data); i += pf.FragmentSize {
		end := i + pf.FragmentSize
		if end > len(data) {
			end = len(data)
		}
		fragments = append(fragments, data[i:end])
	}

	return fragments
}

func (pf *PacketFragmenter) Defragment(fragments [][]byte) []byte {
	var totalLen int
	for _, f := range fragments {
		totalLen += len(f)
	}

	result := make([]byte, 0, totalLen)
	for _, f := range fragments {
		result = append(result, f...)
	}

	return result
}

type TrafficPadder struct {
	MinPad int
	MaxPad int
}

func NewTrafficPadder(minPad, maxPad int) *TrafficPadder {
	if minPad == 0 {
		minPad = 40
	}
	if maxPad == 0 {
		maxPad = 200
	}
	return &TrafficPadder{MinPad: minPad, MaxPad: maxPad}
}

func (tp *TrafficPadder) Pad(data []byte) []byte {
	padSize := tp.MinPad + randInt(tp.MaxPad-tp.MinPad+1)
	padded := make([]byte, len(data)+padSize)
	copy(padded, data)
	randomizeBuffer(padded[len(data):])
	return padded
}

func (tp *TrafficPadder) Unpad(data []byte, originalLen int) []byte {
	if len(data) < originalLen {
		return data
	}
	return data[:originalLen]
}

type TimingRandomizer struct {
	JitterMs int
}

func NewTimingRandomizer(jitterMs int) *TimingRandomizer {
	if jitterMs == 0 {
		jitterMs = 5
	}
	return &TimingRandomizer{JitterMs: jitterMs}
}

func (tr *TimingRandomizer) Randomize() {
	if tr.JitterMs > 0 {
		time.Sleep(time.Duration(randInt(tr.JitterMs*2)) * time.Millisecond)
	}
}
