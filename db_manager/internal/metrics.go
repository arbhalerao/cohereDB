package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "meerkat",
		Name:      "requests_total",
		Help:      "Total number of key-value requests",
	}, []string{"operation", "status"})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "meerkat",
		Name:      "request_duration_seconds",
		Help:      "Duration of key-value requests in seconds",
		Buckets:   prometheus.DefBuckets,
	}, []string{"operation"})

	ActiveServers = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "meerkat",
		Name:      "active_servers",
		Help:      "Number of active database servers in the cluster",
	})

	ReplicationWrites = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "meerkat",
		Name:      "replication_writes_total",
		Help:      "Total replication write attempts",
	}, []string{"status"})

	KeysMigrated = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "meerkat",
		Name:      "keys_migrated_total",
		Help:      "Total number of keys migrated during node add/remove",
	}, []string{"event"})
)
