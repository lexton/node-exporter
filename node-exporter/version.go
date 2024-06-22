package nodeexporter

import (
	"bufio"
	"fmt"
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

var kernelVersionDesc = prometheus.NewDesc(
	"node_kernel_version",
	"Kernel version of the node.",
	[]string{"version"}, nil,
)

func (e *exporter) describeKernelVersionStats(ch chan<- *prometheus.Desc) {
	ch <- kernelVersionDesc
}

func (e *exporter) collectKernelVersion(ch chan<- prometheus.Metric) {
	file, err := os.Open("/proc/version")
	if err != nil {
		fmt.Printf("Error opening /proc/version: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		kernelVersion := scanner.Text()
		ch <- prometheus.MustNewConstMetric(kernelVersionDesc, prometheus.GaugeValue, 1, kernelVersion)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading /proc/version: %v", err)
	}
}
