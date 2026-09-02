package metrics

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsCollector struct {
	mu     sync.RWMutex
	server *http.Server

	activeConnections prometheus.Gauge
	totalConnections  prometheus.Counter
	bytesTransmitted  prometheus.Counter
	bytesReceived     prometheus.Counter
	requestsTotal     prometheus.Counter
	requestDuration   prometheus.Histogram
	errorsTotal       prometheus.Counter
	nodeLatency       *prometheus.GaugeVec
	activeNodes       prometheus.Gauge
	protocolUsage     *prometheus.CounterVec
	activeUsers       prometheus.Gauge
	revenueTotal      prometheus.Counter
}

func NewMetricsCollector(port int) *MetricsCollector {
	m := &MetricsCollector{
		activeConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "nebula_active_connections",
			Help: "Current active VPN connections",
		}),
		totalConnections: promauto.NewCounter(prometheus.CounterOpts{
			Name: "nebula_total_connections_total",
			Help: "Total connections established",
		}),
		bytesTransmitted: promauto.NewCounter(prometheus.CounterOpts{
			Name: "nebula_bytes_transmitted_total",
			Help: "Total bytes transmitted",
		}),
		bytesReceived: promauto.NewCounter(prometheus.CounterOpts{
			Name: "nebula_bytes_received_total",
			Help: "Total bytes received",
		}),
		requestsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "nebula_http_requests_total",
			Help: "Total HTTP requests",
		}),
		requestDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "nebula_http_request_duration_seconds",
			Help:    "HTTP request duration",
			Buckets: prometheus.DefBuckets,
		}),
		errorsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "nebula_errors_total",
			Help: "Total errors",
		}),
		nodeLatency: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Name: "nebula_node_latency_ms",
			Help: "Node latency in milliseconds",
		}, []string{"node", "country"}),
		activeNodes: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "nebula_active_nodes",
			Help: "Number of active VPN nodes",
		}),
		protocolUsage: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "nebula_protocol_usage_total",
			Help: "Protocol usage count",
		}, []string{"protocol"}),
		activeUsers: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "nebula_active_users",
			Help: "Currently active users",
		}),
		revenueTotal: promauto.NewCounter(prometheus.CounterOpts{
			Name: "nebula_revenue_total_cents",
			Help: "Total revenue in cents",
		}),
	}

	m.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: promhttp.Handler(),
	}

	return m
}

func (m *MetricsCollector) Start(ctx context.Context) {
	go func() {
		<-ctx.Done()
		m.server.Shutdown(context.Background())
	}()

	go m.server.ListenAndServe()
}

func (m *MetricsCollector) IncConnections() { m.totalConnections.Inc() }
func (m *MetricsCollector) DecConnections() { m.activeConnections.Dec() }
func (m *MetricsCollector) SetActiveConnections(v int) { m.activeConnections.Set(float64(v)) }
func (m *MetricsCollector) AddBytesTX(v int64) { m.bytesTransmitted.Add(float64(v)) }
func (m *MetricsCollector) AddBytesRX(v int64) { m.bytesReceived.Add(float64(v)) }
func (m *MetricsCollector) IncRequests() { m.requestsTotal.Inc() }
func (m *MetricsCollector) ObserveRequestDuration(d time.Duration) { m.requestDuration.Observe(d.Seconds()) }
func (m *MetricsCollector) IncErrors() { m.errorsTotal.Inc() }
func (m *MetricsCollector) SetNodeLatency(nodeID, country string, ms float64) { m.nodeLatency.WithLabelValues(nodeID, country).Set(ms) }
func (m *MetricsCollector) SetActiveNodes(v int) { m.activeNodes.Set(float64(v)) }
func (m *MetricsCollector) IncProtocol(protocol string) { m.protocolUsage.WithLabelValues(protocol).Inc() }
func (m *MetricsCollector) SetActiveUsers(v int) { m.activeUsers.Set(float64(v)) }
func (m *MetricsCollector) AddRevenue(cents int64) { m.revenueTotal.Add(float64(cents)) }
