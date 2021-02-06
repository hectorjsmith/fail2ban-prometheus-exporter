package main

import (
	"fail2ban-prometheus-exporter/cfg"
	fail2banDb "fail2ban-prometheus-exporter/db"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

const namespace = "fail2ban"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"

	metricUp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Was the last fail2ban query successful.",
		nil, nil,
	)
	metricBannedIpsPerJail = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "banned_ips"),
		"Number of banned IPs stored in the database (per jail).",
		[]string{"jail"}, nil,
	)
	metricBadIpsPerJail = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "bad_ips"),
		"Number of bad IPs stored in the database (per jail).",
		[]string{"jail"}, nil,
	)
)

type Exporter struct {
	db *fail2banDb.Fail2BanDB
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- metricUp
	ch <- metricBadIpsPerJail
	ch <- metricBannedIpsPerJail
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		metricUp, prometheus.GaugeValue, 1,
	)
	e.collectBadIpsPerJailMetrics(ch)
	e.collectBannedIpsPerJailMetrics(ch)
}

func (e *Exporter) collectBadIpsPerJailMetrics(ch chan<- prometheus.Metric) {
	jailNameToCountMap, err := e.db.CountBadIpsPerJail()
	if err != nil {
		log.Print(err)
	}

	for jailName, count := range jailNameToCountMap {
		ch <- prometheus.MustNewConstMetric(
			metricBadIpsPerJail, prometheus.GaugeValue, float64(count), jailName,
		)
	}
}

func (e *Exporter) collectBannedIpsPerJailMetrics(ch chan<- prometheus.Metric) {
	jailNameToCountMap, err := e.db.CountBannedIpsPerJail()
	if err != nil {
		log.Print(err)
	}

	for jailName, count := range jailNameToCountMap {
		ch <- prometheus.MustNewConstMetric(
			metricBannedIpsPerJail, prometheus.GaugeValue, float64(count), jailName,
		)
	}
}

func printAppVersion() {
	fmt.Println(version)
	fmt.Printf("    build date:  %s\r\n    commit hash: %s\r\n    built by:    %s\r\n", date, commit, builtBy)
}

func main() {
	appSettings := cfg.Parse()
	if appSettings.VersionMode {
		printAppVersion()
	} else {
		log.Print("starting fail2ban exporter")

		exporter := &Exporter{
			db: fail2banDb.MustConnectToDb(appSettings.Fail2BanDbPath),
		}
		prometheus.MustRegister(exporter)

		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appSettings.MetricsPort), nil))
	}
}
