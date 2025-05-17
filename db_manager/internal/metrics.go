package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "coheredb",
		Name:      "requests_total",
		Help:      "Total number of key-value requests",
	}, []string{"operation", "status"})

	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "coheredb",
		Name:      "request_duration_seconds",
		Help:      "Duration of key-value requests in seconds",
		Buckets:   prometheus.DefBuckets,
	}, []string{"operation"})

	ActiveServers = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "coheredb",
		Name:      "active_servers",
		Help:      "Number of active database servers in the cluster",
	})

	ReplicationWrites = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "coheredb",
		Name:      "replication_writes_total",
		Help:      "Total replication write attempts",
	}, []string{"status"})

	KeysMigrated = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "coheredb",
		Name:      "keys_migrated_total",
		Help:      "Total number of keys migrated during node add/remove",
	}, []string{"event"})
)
