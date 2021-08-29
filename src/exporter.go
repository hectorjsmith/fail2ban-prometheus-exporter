package main

import (
	"fail2ban-prometheus-exporter/cfg"
	fail2banDb "fail2ban-prometheus-exporter/db"
	"fail2ban-prometheus-exporter/socket"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace     = "fail2ban"
	sockNamespace = "f2b"
)

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
	metricEnabledJails = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "enabled_jails"),
		"Enabled jails.",
		[]string{"jail"}, nil,
	)
	metricErrorCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "errors"),
		"Number of errors found since startup.",
		[]string{"type"}, nil,
	)
	metricServerPing = prometheus.NewDesc(
		prometheus.BuildFQName(sockNamespace, "", "up"),
		"Check if the fail2ban server is up",
		nil, nil,
	)
)

type Exporter struct {
	db           *fail2banDb.Fail2BanDB
	socket       *socket.Fail2BanSocket
	lastError    error
	dbErrorCount int
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	if e.db != nil {
		ch <- metricUp
		ch <- metricBadIpsPerJail
		ch <- metricBannedIpsPerJail
		ch <- metricEnabledJails
		ch <- metricErrorCount
	}
	if e.socket != nil {
		ch <- metricServerPing
	}
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if e.db != nil {
		e.collectBadIpsPerJailMetrics(ch)
		e.collectBannedIpsPerJailMetrics(ch)
		e.collectEnabledJailMetrics(ch)
		e.collectUpMetric(ch)
		e.collectErrorCountMetric(ch)
	}
	if e.socket != nil {
		e.collectServerPingMetric(ch)
	}
}

func (e *Exporter) collectUpMetric(ch chan<- prometheus.Metric) {
	var upMetricValue float64 = 1
	if e.lastError != nil {
		upMetricValue = 0
	}
	ch <- prometheus.MustNewConstMetric(
		metricUp, prometheus.GaugeValue, upMetricValue,
	)
}

func (e *Exporter) collectErrorCountMetric(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		metricErrorCount, prometheus.CounterValue, float64(e.dbErrorCount), "db",
	)
}

func (e *Exporter) collectBadIpsPerJailMetrics(ch chan<- prometheus.Metric) {
	jailNameToCountMap, err := e.db.CountBadIpsPerJail()
	e.lastError = err

	if err != nil {
		e.dbErrorCount++
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
	e.lastError = err

	if err != nil {
		e.dbErrorCount++
		log.Print(err)
	}

	for jailName, count := range jailNameToCountMap {
		ch <- prometheus.MustNewConstMetric(
			metricBannedIpsPerJail, prometheus.GaugeValue, float64(count), jailName,
		)
	}
}

func (e *Exporter) collectEnabledJailMetrics(ch chan<- prometheus.Metric) {
	jailNameToEnabledMap, err := e.db.JailNameToEnabledValue()
	e.lastError = err

	if err != nil {
		e.dbErrorCount++
		log.Print(err)
	}

	for jailName, count := range jailNameToEnabledMap {
		ch <- prometheus.MustNewConstMetric(
			metricEnabledJails, prometheus.GaugeValue, float64(count), jailName,
		)
	}
}

func (e *Exporter) collectServerPingMetric(ch chan<- prometheus.Metric) {
	pingSuccess := e.socket.Ping()
	var pingSuccessInt float64 = 1
	if !pingSuccess {
		pingSuccessInt = 0
	}
	ch <- prometheus.MustNewConstMetric(
		metricServerPing, prometheus.GaugeValue, pingSuccessInt,
	)
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

		exporter := &Exporter{}
		if appSettings.Fail2BanDbPath != "" {
			exporter.db = fail2banDb.MustConnectToDb(appSettings.Fail2BanDbPath)
		}
		if appSettings.Fail2BanSocketPath != "" {
			exporter.socket = socket.MustConnectToSocket(appSettings.Fail2BanSocketPath)
		}
		prometheus.MustRegister(exporter)

		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appSettings.MetricsPort), nil))
	}
}
