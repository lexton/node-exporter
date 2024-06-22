package nodeexporter

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	diskReadsCompletedDesc = prometheus.NewDesc("node_disk_reads_completed_total", "Total number of reads completed successfully.", []string{"device"}, nil)
	diskReadsMergedDesc    = prometheus.NewDesc("node_disk_reads_merged_total", "Total number of reads merged.", []string{"device"}, nil)
	diskSectorsReadDesc    = prometheus.NewDesc("node_disk_sectors_read_total", "Total number of sectors read successfully.", []string{"device"}, nil)
	diskReadTimeDesc       = prometheus.NewDesc("node_disk_read_time_seconds_total", "Total time spent reading from disk (ms).", []string{"device"}, nil)

	diskWritesCompletedDesc = prometheus.NewDesc("node_disk_writes_completed_total", "Total number of writes completed successfully.", []string{"device"}, nil)
	diskWritesMergedDesc    = prometheus.NewDesc("node_disk_writes_merged_total", "Total number of writes merged.", []string{"device"}, nil)
	diskSectorsWrittenDesc  = prometheus.NewDesc("node_disk_sectors_written_total", "Total number of sectors written successfully.", []string{"device"}, nil)
	diskWriteTimeDesc       = prometheus.NewDesc("node_disk_write_time_seconds_total", "Total time spent writing to disk (ms).", []string{"device"}, nil)
)

func (e *exporter) describeDiskStats(ch chan<- *prometheus.Desc) {
	ch <- diskReadsCompletedDesc
	ch <- diskReadsMergedDesc
	ch <- diskSectorsReadDesc
	ch <- diskReadTimeDesc
	ch <- diskWritesCompletedDesc
	ch <- diskWritesMergedDesc
	ch <- diskSectorsWrittenDesc
	ch <- diskWriteTimeDesc
}

func (e *exporter) collectDiskStats(ch chan<- prometheus.Metric) {
	file, err := os.Open("/proc/diskstats")
	if err != nil {
		fmt.Printf("Error opening /proc/diskstats: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) < 14 {
			continue
		}

		device := parts[2]

		// See: https://www.kernel.org/doc/Documentation/admin-guide/iostats.rst
		// This is going to depend on the kernel version you are running
		// We assume version 4.18+

		// TODO(lexton): Add discard/flush counters here
		// TODO(lexton): Add active io timers here - needs more

		// NOTE: add 2 to the index values in the docs
		ch <- prometheus.MustNewConstMetric(diskReadsCompletedDesc, prometheus.CounterValue, parseFloat(parts[3]), device)
		ch <- prometheus.MustNewConstMetric(diskReadsMergedDesc, prometheus.CounterValue, parseFloat(parts[4]), device)
		ch <- prometheus.MustNewConstMetric(diskSectorsReadDesc, prometheus.CounterValue, parseFloat(parts[5]), device)
		ch <- prometheus.MustNewConstMetric(diskReadTimeDesc, prometheus.CounterValue, parseFloat(parts[6])/1000, device)

		ch <- prometheus.MustNewConstMetric(diskWritesCompletedDesc, prometheus.CounterValue, parseFloat(parts[7]), device)
		ch <- prometheus.MustNewConstMetric(diskWritesMergedDesc, prometheus.CounterValue, parseFloat(parts[8]), device)
		ch <- prometheus.MustNewConstMetric(diskSectorsWrittenDesc, prometheus.CounterValue, parseFloat(parts[9]), device)
		ch <- prometheus.MustNewConstMetric(diskWriteTimeDesc, prometheus.CounterValue, parseFloat(parts[10])/1000, device)

	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading /proc/diskstat: %v", err)
	}
}
