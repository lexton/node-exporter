package nodeexporter

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var collectDuration = prometheus.NewDesc("node_metrics_collect_duration_seconds", "Duration of the collect call", nil, nil)

type exporter struct{}

func New() prometheus.Collector {
	return &exporter{}
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	e.describeCPUStats(ch)
	e.describeDiskStats(ch)
	e.describeUptimeStats(ch)
	e.describeLoadAvgStats(ch)
	e.describeNICStats(ch)
	e.describeMemStats(ch)
	e.describeKernelVersionStats(ch)
	ch <- collectDuration
}

// Collect is called sync w/ metrics collections tasks
// this is important to reduce metric propegation delays when
// compared to a background goroutine that would update metrics
// it's important to ensure that collection time remains under 10s
func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	e.collectDiskStats(ch)
	e.collectUptime(ch)
	e.collectLoadAvg(ch)
	e.collectNICStats(ch)
	e.collectMemStats(ch)
	e.collectKernelVersion(ch)

	ch <- prometheus.MustNewConstMetric(collectDuration, prometheus.CounterValue, time.Since(start).Seconds())
}
