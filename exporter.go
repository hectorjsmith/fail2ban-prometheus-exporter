package main

import (
	fail2banDb "fail2ban-prometheus-exporter/db"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

const namespace = "fail2ban"

var (
	db       = fail2banDb.MustConnectToDb("fail2ban.sqlite3")
	metricUp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last fail2ban query successful.",
		nil, nil,
	)
	metricBadIpTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "badip_total"),
		"Total number of bad IPs stored in the database.",
		nil, nil,
	)
)

type Exporter struct {
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- metricUp
	ch <- metricBadIpTotal
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		metricUp, prometheus.GaugeValue, 1,
	)
	ch <- *collectTotalBadIpMetric()
}

func collectTotalBadIpMetric() *prometheus.Metric {
	count, _ := db.CountTotalBadIps()
	metric := prometheus.MustNewConstMetric(
		metricBadIpTotal, prometheus.GaugeValue, float64(count),
	)
	return &metric
}

func main() {
	log.Print("starting fail2ban exporter")

	exporter := &Exporter{}
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
