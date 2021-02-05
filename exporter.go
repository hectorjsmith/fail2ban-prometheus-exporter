package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

const namespace = "fail2ban"

var up = prometheus.NewDesc(
	prometheus.BuildFQName(namespace, "", "up"),
	"Was the last fail2ban query successful.",
	nil, nil,
)

type Exporter struct {
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		up, prometheus.GaugeValue, 1,
	)
}

func main() {
	log.Print("starting fail2ban exporter")

	exporter := &Exporter{}
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
