package node

import (
	"math"
	"sort"
	"sync"

	"github.com/nebula-vpn/backend/internal/db"
)

type GeoBalancer struct {
	mu      sync.RWMutex
	servers []*db.VPNServer
}

func NewGeoBalancer() *GeoBalancer {
	return &GeoBalancer{}
}

func (g *GeoBalancer) UpdateServers(servers []*db.VPNServer) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.servers = servers
}

func (g *GeoBalancer) SelectOptimal(clientLat float64, clientLon float64, protocol string, country string) *db.VPNServer {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if len(g.servers) == 0 {
		return nil
	}

	type scored struct {
		server *db.VPNServer
		score  float64
	}

	candidates := make([]scored, 0)

	for _, s := range g.servers {
		if !s.IsActive || !s.IsHealthy {
			continue
		}
		if country != "" && s.Country != country {
			continue
		}

		dist := haversine(clientLat, clientLon, s.Latitude, s.Longitude)

		loadScore := 1.0 - float64(s.LoadPercent)/100.0
		distScore := 1.0 / (1.0 + dist/1000.0)
		capacityScore := float64(s.MaxConnections-s.CurrentConnections) / float64(s.MaxConnections)

		score := 0.4*distScore + 0.35*loadScore + 0.25*capacityScore

		candidates = append(candidates, scored{server: s, score: score})
	}

	if len(candidates) == 0 {
		return nil
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].score > candidates[j].score
	})

	return candidates[0].server
}

func (g *GeoBalancer) SelectForFailover(excludeID string) *db.VPNServer {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, s := range g.servers {
		if s.ID.String() != excludeID && s.IsActive && s.IsHealthy {
			return s
		}
	}
	return nil
}

func (g *GeoBalancer) HealthyCount() int {
	g.mu.RLock()
	defer g.mu.RUnlock()

	count := 0
	for _, s := range g.servers {
		if s.IsActive && s.IsHealthy {
			count++
		}
	}
	return count
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0

	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180.0)*math.Cos(lat2*math.Pi/180.0)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
