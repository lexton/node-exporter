package nodeexporter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	memTotalDesc     = prometheus.NewDesc("node_mem_total_bytes", "Total memory in bytes", nil, nil)
	memFreeDesc      = prometheus.NewDesc("node_mem_free_bytes", "Free memory in bytes", nil, nil)
	memAvailableDesc = prometheus.NewDesc("node_mem_available_bytes", "Available memory in bytes", nil, nil)
	swapTotalDesc    = prometheus.NewDesc("node_mem_swap_total_bytes", "Total swap space in bytes", nil, nil)
	swapFreeDesc     = prometheus.NewDesc("node_mem_swap_free_bytes", "Free swap space in bytes", nil, nil)
	mappedDesc       = prometheus.NewDesc("node_mem_mapped_bytes", "Mapped memory in bytes", nil, nil)
	shmemDesc        = prometheus.NewDesc("node_mem_shmem_bytes", "Shared memory in bytes", nil, nil)
)

func (e *exporter) describeMemStats(ch chan<- *prometheus.Desc) {
	ch <- memTotalDesc
	ch <- memFreeDesc
	ch <- memAvailableDesc
	ch <- swapTotalDesc
	ch <- swapFreeDesc
	ch <- mappedDesc
	ch <- shmemDesc
}

func (e *exporter) collectMemStats(ch chan<- prometheus.Metric) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Println("Error opening /proc/meminfo:", err)
		return
	}
	defer file.Close()

	memStats := map[string]float64{
		"MemTotal":     0,
		"MemFree":      0,
		"MemAvailable": 0,
		"SwapTotal":    0,
		"SwapFree":     0,
		"Mapped":       0,
		"Shmem":        0,
	}

	converter := map[string]float64{
		"kB": 1024,
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSuffix(fields[0], ":")
		if _, ok := memStats[key]; ok {
			memStats[key] = parseFloat(fields[1]) * converter[fields[2]]
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading /proc/meminfo:", err)
		return
	}

	ch <- prometheus.MustNewConstMetric(memTotalDesc, prometheus.GaugeValue, memStats["MemTotal"])
	ch <- prometheus.MustNewConstMetric(memFreeDesc, prometheus.GaugeValue, memStats["MemFree"])
	ch <- prometheus.MustNewConstMetric(memAvailableDesc, prometheus.GaugeValue, memStats["MemAvailable"])
	ch <- prometheus.MustNewConstMetric(swapTotalDesc, prometheus.GaugeValue, memStats["SwapTotal"])
	ch <- prometheus.MustNewConstMetric(swapFreeDesc, prometheus.GaugeValue, memStats["SwapFree"])
	ch <- prometheus.MustNewConstMetric(mappedDesc, prometheus.GaugeValue, memStats["Mapped"])
	ch <- prometheus.MustNewConstMetric(shmemDesc, prometheus.GaugeValue, memStats["Shmem"])
}
