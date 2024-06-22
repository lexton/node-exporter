package main

import (
	"fmt"
	"net/http"
	"os"

	nodeexporter "github.com/lexton/node-exporter/v2/node-exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(nodeexporter.New())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	fmt.Println("Starting Node Exporter :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Failed to start server", err)
	}
}
