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
	cpuDesc          = prometheus.NewDesc("node_cpu_seconds_total", "Seconds the CPUs spent in each mode.", []string{"cpu", "mode"}, nil)
	intrDesc         = prometheus.NewDesc("node_intr_total", "Total number of interrupts serviced.", nil, nil)
	ctxtDesc         = prometheus.NewDesc("node_context_switches_total", "Total number of context switches.", nil, nil)
	btimeDesc        = prometheus.NewDesc("node_boot_time_seconds", "Node boot time, in seconds since epoch.", nil, nil)
	processesDesc    = prometheus.NewDesc("node_processes_total", "Total number of processes.", nil, nil)
	procsRunningDesc = prometheus.NewDesc("node_procs_running", "Number of processes in runnable state.", nil, nil)
	procsBlockedDesc = prometheus.NewDesc("node_procs_blocked", "Number of processes blocked.", nil, nil)
	softirqDesc      = prometheus.NewDesc("node_softirq_total", "Total number of soft IRQs.", []string{"type"}, nil)
)

func (e *exporter) describeCPUStats(ch chan<- *prometheus.Desc) {
	ch <- cpuDesc
	ch <- intrDesc
	ch <- ctxtDesc
	ch <- btimeDesc
	ch <- processesDesc
	ch <- procsRunningDesc
	ch <- procsBlockedDesc
	ch <- softirqDesc
}

func (e *exporter) collectCPUStats(ch chan<- prometheus.Metric) {
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

			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, parseFloat(parts[1])/100, cpu, "user")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, parseFloat(parts[2])/100, cpu, "nice")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, parseFloat(parts[3])/100, cpu, "system")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, parseFloat(parts[4])/100, cpu, "idle")

			// The CPU will not wait for I/O to complete; iowait is the time that a task is waiting for I/O to complete.
			// When a CPU goes into idle state for outstanding task I/O, another task will be scheduled on this CPU.
			//ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, parseFloat(parts[5])/100, cpu, "iowait")

			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, parseFloat(parts[6])/100, cpu, "irq")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, parseFloat(parts[7])/100, cpu, "softirq")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, parseFloat(parts[8])/100, cpu, "steal")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, parseFloat(parts[9])/100, cpu, "guest")
			ch <- prometheus.MustNewConstMetric(cpuDesc, prometheus.CounterValue, parseFloat(parts[10])/100, cpu, "guest_nice")

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
