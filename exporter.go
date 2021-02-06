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
	metricBadIpsPerJail = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "bad_ips"),
		"Number of bad IPs stored in the database (per jail).",
		[]string{"jail"}, nil,
	)
)

type Exporter struct {
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- metricUp
	ch <- metricBadIpsPerJail
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		metricUp, prometheus.GaugeValue, 1,
	)
	collectBadIpsPerJailMetrics(ch)
}

func collectBadIpsPerJailMetrics(ch chan<- prometheus.Metric) {
	jailNameToCountMap, err := db.CountBadIpsPerJail()
	if err != nil {
		log.Print(err)
	}

	for jailName, count := range jailNameToCountMap {
		ch <- prometheus.MustNewConstMetric(
			metricBadIpsPerJail, prometheus.GaugeValue, float64(count), jailName,
		)
	}
}

func main() {
	log.Print("starting fail2ban exporter")

	exporter := &Exporter{}
	prometheus.MustRegister(exporter)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}
