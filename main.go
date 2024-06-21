package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cpuDesc = prometheus.NewDesc(
		"node_cpu_seconds_total",
		"Seconds the CPUs spent in each mode.",
		[]string{"cpu", "mode"}, nil,
	)
	intrDesc = prometheus.NewDesc(
		"node_intr_total",
		"Total number of interrupts serviced.",
		nil, nil,
	)
	ctxtDesc = prometheus.NewDesc(
		"node_context_switches_total",
		"Total number of context switches.",
		nil, nil,
	)
	btimeDesc = prometheus.NewDesc(
		"node_boot_time_seconds",
		"Node boot time, in seconds since epoch.",
		nil, nil,
	)
	processesDesc = prometheus.NewDesc(
		"node_processes_total",
		"Total number of processes.",
		nil, nil,
	)
	procsRunningDesc = prometheus.NewDesc(
		"node_procs_running",
		"Number of processes in runnable state.",
		nil, nil,
	)
	procsBlockedDesc = prometheus.NewDesc(
		"node_procs_blocked",
		"Number of processes blocked.",
		nil, nil,
	)
	softirqDesc = prometheus.NewDesc(
		"node_softirq_total",
		"Total number of soft IRQs.",
		[]string{"type"}, nil,
	)
)

type exporter struct{}

func newExporter() *exporter {
	return &exporter{}
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- cpuDesc
	ch <- intrDesc
	ch <- ctxtDesc
	ch <- btimeDesc
	ch <- processesDesc
	ch <- procsRunningDesc
	ch <- procsBlockedDesc
	ch <- softirqDesc
}

// Collect is called sync w/ metrics collections tasks
// this is important to reduce metric propegation delays when
// compared to a background goroutine that would update metrics
// it's important to ensure that collection time remains under 10s
func (e *exporter) Collect(ch chan<- prometheus.Metric) {
	log.Println("Collecting Metrics")
	// See man 5 proc
	file, err := os.Open("/proc/stat")
	if err != nil {
		fmt.Printf("Error opening /proc/stat: %v", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		switch {
		case strings.HasPrefix(line, "cpu"):
			if len(parts) < 11 {
				continue
			}
			cpu := parts[0]
			user, _ := strconv.ParseFloat(parts[1], 64)
			nice, _ := strconv.ParseFloat(parts[2], 64)
			system, _ := strconv.ParseFloat(parts[3], 64)
			idle, _ := strconv.ParseFloat(parts[4], 64)
			// The CPU will not wait for I/O to complete; iowait is the time that a task is waiting for I/O to complete.
			// When a CPU goes into idle state for outstanding task I/O, another task will be scheduled on this CPU.
			// iowait, _ := strconv.ParseFloat(parts[5], 64)
			irq, _ := strconv.ParseFloat(parts[6], 64)
			softirq, _ := strconv.ParseFloat(parts[7], 64)
			steal, _ := strconv.ParseFloat(parts[8], 64)
			guest, _ := strconv.ParseFloat(parts[9], 64)
			guestNice, _ := strconv.ParseFloat(parts[10], 64)

			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, user/100, cpu, "user")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, nice/100, cpu, "nice")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, system/100, cpu, "system")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, idle/100, cpu, "idle")
			//ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, iowait/100, cpu, "iowait")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, irq/100, cpu, "irq")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, softirq/100, cpu, "softirq")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, steal/100, cpu, "steal")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, guest/100, cpu, "guest")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, guestNice/100, cpu, "guest_nice")

		case strings.HasPrefix(line, "intr"):
			intr, _ := strconv.ParseFloat(parts[1], 64)
			ch <- prometheus.MustNewConstMetric(intrDesc, prometheus.CounterValue, intr)

		case strings.HasPrefix(line, "ctxt"):
			ctxt, _ := strconv.ParseFloat(parts[1], 64)
			ch <- prometheus.MustNewConstMetric(ctxtDesc, prometheus.CounterValue, ctxt)

		case strings.HasPrefix(line, "btime"):
			btime, _ := strconv.ParseFloat(parts[1], 64)
			ch <- prometheus.MustNewConstMetric(btimeDesc, prometheus.GaugeValue, btime)

		case strings.HasPrefix(line, "processes"):
			processes, _ := strconv.ParseFloat(parts[1], 64)
			ch <- prometheus.MustNewConstMetric(processesDesc, prometheus.CounterValue, processes)

		case strings.HasPrefix(line, "procs_running"):
			procsRunning, _ := strconv.ParseFloat(parts[1], 64)
			ch <- prometheus.MustNewConstMetric(procsRunningDesc, prometheus.GaugeValue, procsRunning)

		case strings.HasPrefix(line, "procs_blocked"):
			procsBlocked, _ := strconv.ParseFloat(parts[1], 64)
			ch <- prometheus.MustNewConstMetric(procsBlockedDesc, prometheus.GaugeValue, procsBlocked)

		case strings.HasPrefix(line, "softirq"):
			for i, irq := range parts[1:] {
				value, _ := strconv.ParseFloat(irq, 64)
				ch <- prometheus.MustNewConstMetric(softirqDesc, prometheus.CounterValue, value, fmt.Sprintf("%d", i))
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading /proc/stat: %v", err)
	}
}

func main() {
	reg := prometheus.NewRegistry()
	exporter := newExporter()
	reg.MustRegister(exporter)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	fmt.Println("Starting Node Exporter :" + port)
	http.ListenAndServe(":"+port, nil)
}
