package analytics

import (
	"sync"
	"time"
)

type AnalyticsCollector struct {
	mu sync.RWMutex

	connectionsByCountry map[string]int64
	connectionsByProtocol map[string]int64
	dailyActiveUsers     map[string]int64
	revenueByDay         map[string]int64
	registrationCount    int64
}

func NewAnalyticsCollector() *AnalyticsCollector {
	return &AnalyticsCollector{
		connectionsByCountry: make(map[string]int64),
		connectionsByProtocol: make(map[string]int64),
		dailyActiveUsers:     make(map[string]int64),
		revenueByDay:         make(map[string]int64),
	}
}

func (ac *AnalyticsCollector) RecordConnection(country, protocol string) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.connectionsByCountry[country]++
	ac.connectionsByProtocol[protocol]++
}

func (ac *AnalyticsCollector) RecordActiveUser(userID string) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	today := time.Now().UTC().Format("2006-01-02")
	// Simple counting, dedup should use Redis in production
	ac.dailyActiveUsers[today]++
}

func (ac *AnalyticsCollector) RecordRevenue(amount int64) {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	today := time.Now().UTC().Format("2006-01-02")
	ac.revenueByDay[today] += amount
}

func (ac *AnalyticsCollector) RecordRegistration() {
	ac.mu.Lock()
	defer ac.mu.Unlock()
	ac.registrationCount++
}

func (ac *AnalyticsCollector) GetReport() *AnalyticsReport {
	ac.mu.RLock()
	defer ac.mu.RUnlock()

	return &AnalyticsReport{
		ConnectionsByCountry: cloneMap(ac.connectionsByCountry),
		ConnectionsByProtocol: cloneMap(ac.connectionsByProtocol),
		DailyActiveUsers:     cloneMap(ac.dailyActiveUsers),
		RevenueByDay:         cloneMap(ac.revenueByDay),
		TotalRegistrations:   ac.registrationCount,
	}
}

type AnalyticsReport struct {
	ConnectionsByCountry  map[string]int64 `json:"connections_by_country"`
	ConnectionsByProtocol map[string]int64 `json:"connections_by_protocol"`
	DailyActiveUsers      map[string]int64 `json:"daily_active_users"`
	RevenueByDay          map[string]int64 `json:"revenue_by_day"`
	TotalRegistrations    int64            `json:"total_registrations"`
}

func cloneMap(src map[string]int64) map[string]int64 {
	dst := make(map[string]int64, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
