package nodeexporter

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	load1Desc = prometheus.NewDesc(
		"node_load1",
		"1 minute load average.",
		nil, nil,
	)
	load5Desc = prometheus.NewDesc(
		"node_load5",
		"5 minute load average.",
		nil, nil,
	)
	load15Desc = prometheus.NewDesc(
		"node_load15",
		"15 minute load average.",
		nil, nil,
	)
)

func (e *exporter) describeLoadAvgStats(ch chan<- *prometheus.Desc) {
	ch <- load1Desc
	ch <- load5Desc
	ch <- load15Desc
}

func (e *exporter) collectLoadAvg(ch chan<- prometheus.Metric) {
	// Collect load average
	loadavgFile, err := os.Open("/proc/loadavg")
	if err != nil {
		fmt.Printf("Error opening /proc/loadavg: %v", err)
		return
	}
	defer loadavgFile.Close()

	var load1, load5, load15 float64
	if _, err := fmt.Fscanf(loadavgFile, "%f %f %f", &load1, &load5, &load15); err != nil {
		fmt.Printf("Error reading /proc/loadavg: %v", err)
		return
	}
	ch <- prometheus.MustNewConstMetric(load1Desc, prometheus.GaugeValue, load1)
	ch <- prometheus.MustNewConstMetric(load5Desc, prometheus.GaugeValue, load5)
	ch <- prometheus.MustNewConstMetric(load15Desc, prometheus.GaugeValue, load15)
}
