package nodeexporter

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

var uptimeDesc = prometheus.NewDesc(
	"node_uptime_seconds",
	"System uptime in seconds.",
	nil, nil,
)

func (e *exporter) describeUptimeStats(ch chan<- *prometheus.Desc) {
	ch <- uptimeDesc
}

func (e *exporter) collectUptime(ch chan<- prometheus.Metric) {
	// Collect uptime
	uptimeFile, err := os.Open("/proc/uptime")
	if err != nil {
		fmt.Printf("Error opening /proc/uptime: %v", err)
		return
	}
	defer uptimeFile.Close()

	var uptime float64
	if _, err := fmt.Fscanf(uptimeFile, "%f", &uptime); err != nil {
		fmt.Printf("Error reading /proc/uptime: %v", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(uptimeDesc, prometheus.GaugeValue, uptime)

}
