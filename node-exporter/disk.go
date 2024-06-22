package nodeexporter

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	diskReadsCompletedDesc = prometheus.NewDesc(
		"node_disk_reads_completed_total",
		"Total number of reads completed successfully.",
		[]string{"device"}, nil,
	)
	diskReadsMergedDesc = prometheus.NewDesc(
		"node_disk_reads_merged_total",
		"Total number of reads merged.",
		[]string{"device"}, nil,
	)
	diskSectorsReadDesc = prometheus.NewDesc(
		"node_disk_sectors_read_total",
		"Total number of sectors read successfully.",
		[]string{"device"}, nil,
	)
	diskReadTimeDesc = prometheus.NewDesc(
		"node_disk_read_time_seconds_total",
		"Total time spent reading from disk (ms).",
		[]string{"device"}, nil,
	)
	diskWritesCompletedDesc = prometheus.NewDesc(
		"node_disk_writes_completed_total",
		"Total number of writes completed successfully.",
		[]string{"device"}, nil,
	)
	diskWritesMergedDesc = prometheus.NewDesc(
		"node_disk_writes_merged_total",
		"Total number of writes merged.",
		[]string{"device"}, nil,
	)
	diskSectorsWrittenDesc = prometheus.NewDesc(
		"node_disk_sectors_written_total",
		"Total number of sectors written successfully.",
		[]string{"device"}, nil,
	)
	diskWriteTimeDesc = prometheus.NewDesc(
		"node_disk_write_time_seconds_total",
		"Total time spent writing to disk (ms).",
		[]string{"device"}, nil,
	)
	diskIOTimeDesc = prometheus.NewDesc(
		"node_disk_io_time_seconds_total",
		"Total time spent doing I/Os (ms).",
		[]string{"device"}, nil,
	)
	diskWeightedIOTimeDesc = prometheus.NewDesc(
		"node_disk_weighted_io_time_seconds_total",
		"Total weighted time spent doing I/Os (ms).",
		[]string{"device"}, nil,
	)
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
	ch <- diskIOTimeDesc
	ch <- diskWeightedIOTimeDesc
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
		readsCompleted, _ := strconv.ParseFloat(parts[3], 64)
		readsMerged, _ := strconv.ParseFloat(parts[4], 64)
		sectorsRead, _ := strconv.ParseFloat(parts[5], 64)
		readTime, _ := strconv.ParseFloat(parts[6], 64)
		writesCompleted, _ := strconv.ParseFloat(parts[7], 64)
		writesMerged, _ := strconv.ParseFloat(parts[8], 64)
		sectorsWritten, _ := strconv.ParseFloat(parts[9], 64)
		writeTime, _ := strconv.ParseFloat(parts[10], 64)
		ioTime, _ := strconv.ParseFloat(parts[12], 64)
		weightedIOTime, _ := strconv.ParseFloat(parts[13], 64)

		ch <- prometheus.MustNewConstMetric(diskReadsCompletedDesc, prometheus.CounterValue, readsCompleted, device)
		ch <- prometheus.MustNewConstMetric(diskReadsMergedDesc, prometheus.CounterValue, readsMerged, device)
		ch <- prometheus.MustNewConstMetric(diskSectorsReadDesc, prometheus.CounterValue, sectorsRead, device)
		ch <- prometheus.MustNewConstMetric(diskReadTimeDesc, prometheus.CounterValue, readTime/1000, device)
		ch <- prometheus.MustNewConstMetric(diskWritesCompletedDesc, prometheus.CounterValue, writesCompleted, device)
		ch <- prometheus.MustNewConstMetric(diskWritesMergedDesc, prometheus.CounterValue, writesMerged, device)
		ch <- prometheus.MustNewConstMetric(diskSectorsWrittenDesc, prometheus.CounterValue, sectorsWritten, device)
		ch <- prometheus.MustNewConstMetric(diskWriteTimeDesc, prometheus.CounterValue, writeTime/1000, device)
		ch <- prometheus.MustNewConstMetric(diskIOTimeDesc, prometheus.CounterValue, ioTime/1000, device)
		ch <- prometheus.MustNewConstMetric(diskWeightedIOTimeDesc, prometheus.CounterValue, weightedIOTime/1000, device)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading /proc/diskstat: %v", err)
	}
}
