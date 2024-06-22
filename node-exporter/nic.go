package nodeexporter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	nicUpDesc = prometheus.NewDesc("nic_up", "NIC up status", []string{"nic"}, nil)

	nicRxBytesDesc = prometheus.NewDesc("nic_rx_bytes_total", "Total received bytes", []string{"nic"}, nil)
	nicRxErrDesc   = prometheus.NewDesc("nic_rx_errs_total", "Total received errors", []string{"nic"}, nil)
	nicRxDropDesc  = prometheus.NewDesc("nic_rx_drops_total", "Total received drops", []string{"nic"}, nil)

	nicTxBytesDesc = prometheus.NewDesc("nic_tx_bytes_total", "Total transmitted bytes", []string{"nic"}, nil)
	nicTxErrDesc   = prometheus.NewDesc("nic_tx_errs_total", "Total transmitted errors", []string{"nic"}, nil)
	nicTxDropDesc  = prometheus.NewDesc("nic_tx_drops_total", "Total transmitted drops", []string{"nic"}, nil)
)

func (e *exporter) describeNICStats(ch chan<- *prometheus.Desc) {
	ch <- nicUpDesc
	ch <- nicRxBytesDesc
	ch <- nicRxErrDesc
	ch <- nicRxDropDesc
	ch <- nicTxBytesDesc
	ch <- nicTxErrDesc
	ch <- nicTxDropDesc
}

func (e *exporter) collectNICStats(ch chan<- prometheus.Metric) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		fmt.Println("Error opening /proc/net/dev:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// Skip the first two lines of headers
	scanner.Scan()
	scanner.Scan()

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 17 {
			continue
		}

		nic := strings.TrimSuffix(fields[0], ":")
		up := 1.0 // Assume NIC is up if it has entries in /proc/net/dev

		ch <- prometheus.MustNewConstMetric(nicUpDesc, prometheus.GaugeValue, up, nic)

		ch <- prometheus.MustNewConstMetric(nicRxBytesDesc, prometheus.CounterValue, parseFloat(fields[1]), nic)
		ch <- prometheus.MustNewConstMetric(nicRxErrDesc, prometheus.CounterValue, parseFloat(fields[2]), nic)
		ch <- prometheus.MustNewConstMetric(nicRxDropDesc, prometheus.CounterValue, parseFloat(fields[3]), nic)

		ch <- prometheus.MustNewConstMetric(nicTxBytesDesc, prometheus.CounterValue, parseFloat(fields[9]), nic)
		ch <- prometheus.MustNewConstMetric(nicTxErrDesc, prometheus.CounterValue, parseFloat(fields[10]), nic)
		ch <- prometheus.MustNewConstMetric(nicTxDropDesc, prometheus.CounterValue, parseFloat(fields[11]), nic)

		// TODO(lexton): Add counters for, fifo frame compressed multicast
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading /proc/net/dev:", err)
	}
}
