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

	metricErrorCountNew = prometheus.NewDesc(
		prometheus.BuildFQName(sockNamespace, "", "errors"),
		"Number of errors found since startup",
		[]string{"type"}, nil,
	)
	metricServerUp = prometheus.NewDesc(
		prometheus.BuildFQName(sockNamespace, "", "up"),
		"Check if the fail2ban server is up",
		nil, nil,
	)
	metricJailCount = prometheus.NewDesc(
		prometheus.BuildFQName(sockNamespace, "", "jail_count"),
		"Number of defined jails",
		nil, nil,
	)
	metricJailFailedCurrent = prometheus.NewDesc(
		prometheus.BuildFQName(sockNamespace, "", "jail_failed_current"),
		"Number of current failures on this jail's filter",
		[]string{"jail"}, nil,
	)
	metricJailFailedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(sockNamespace, "", "jail_failed_total"),
		"Number of total failures on this jail's filter",
		[]string{"jail"}, nil,
	)
	metricJailBannedCurrent = prometheus.NewDesc(
		prometheus.BuildFQName(sockNamespace, "", "jail_banned_current"),
		"Number of IPs currently banned in this jail",
		[]string{"jail"}, nil,
	)
	metricJailBannedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(sockNamespace, "", "jail_banned_total"),
		"Total number of IPs banned by this jail (includes expired bans)",
		[]string{"jail"}, nil,
	)
)

type Exporter struct {
	db                         *fail2banDb.Fail2BanDB
	socketPath                 string
	lastError                  error
	dbErrorCount               int
	socketConnectionErrorCount int
	socketRequestErrorCount    int
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	if e.db != nil {
		ch <- metricUp
		ch <- metricBadIpsPerJail
		ch <- metricBannedIpsPerJail
		ch <- metricEnabledJails
		ch <- metricErrorCount
	}
	if e.socketPath != "" {
		ch <- metricServerUp
		ch <- metricJailCount
		ch <- metricJailFailedCurrent
		ch <- metricJailFailedTotal
		ch <- metricJailBannedCurrent
		ch <- metricJailBannedTotal
	}
	ch <- metricErrorCountNew
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if e.db != nil {
		e.collectBadIpsPerJailMetrics(ch)
		e.collectBannedIpsPerJailMetrics(ch)
		e.collectEnabledJailMetrics(ch)
		e.collectUpMetric(ch)
		e.collectErrorCountMetric(ch)
	}
	if e.socketPath != "" {
		s, err := socket.ConnectToSocket(e.socketPath)
		if err != nil {
			log.Printf("error opening socket: %v", err)
			e.socketConnectionErrorCount++
		} else {
			defer s.Close()
		}
		e.collectServerUpMetric(ch, s)
		if err == nil && s != nil {
			e.collectJailMetrics(ch, s)
		}
	}
	e.collectErrorCountMetricNew(ch)
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

func (e *Exporter) collectErrorCountMetricNew(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		metricErrorCountNew, prometheus.CounterValue, float64(e.dbErrorCount), "db",
	)
	ch <- prometheus.MustNewConstMetric(
		metricErrorCountNew, prometheus.CounterValue, float64(e.socketConnectionErrorCount), "socket_conn",
	)
	ch <- prometheus.MustNewConstMetric(
		metricErrorCountNew, prometheus.CounterValue, float64(e.socketRequestErrorCount), "socket_req",
	)
}

func (e *Exporter) collectServerUpMetric(ch chan<- prometheus.Metric, s *socket.Fail2BanSocket) {
	var serverUp float64 = 0
	if s != nil {
		pingSuccess, err := s.Ping()
		if err != nil {
			e.socketRequestErrorCount++
			log.Print(err)
		}
		if err == nil && pingSuccess {
			serverUp = 1
		}
	}
	ch <- prometheus.MustNewConstMetric(
		metricServerUp, prometheus.GaugeValue, serverUp,
	)
}

func (e *Exporter) collectJailMetrics(ch chan<- prometheus.Metric, s *socket.Fail2BanSocket) {
	jails, err := s.GetJails()
	var count float64 = 0
	if err != nil {
		e.socketRequestErrorCount++
		log.Print(err)
	}
	if err == nil {
		count = float64(len(jails))
	}
	ch <- prometheus.MustNewConstMetric(
		metricJailCount, prometheus.GaugeValue, count,
	)

	for i := range jails {
		e.collectJailStatsMetric(ch, s, jails[i])
	}
}

func (e *Exporter) collectJailStatsMetric(ch chan<- prometheus.Metric, s *socket.Fail2BanSocket, jail string) {
	stats, err := s.GetJailStats(jail)
	if err != nil {
		e.socketRequestErrorCount++
		log.Printf("failed to get stats for jail %s: %v", jail, err)
		return
	}

	ch <- prometheus.MustNewConstMetric(
		metricJailFailedCurrent, prometheus.GaugeValue, float64(stats.FailedCurrent), jail,
	)
	ch <- prometheus.MustNewConstMetric(
		metricJailFailedTotal, prometheus.GaugeValue, float64(stats.FailedTotal), jail,
	)
	ch <- prometheus.MustNewConstMetric(
		metricJailBannedCurrent, prometheus.GaugeValue, float64(stats.BannedCurrent), jail,
	)
	ch <- prometheus.MustNewConstMetric(
		metricJailBannedTotal, prometheus.GaugeValue, float64(stats.BannedTotal), jail,
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
			exporter.socketPath = appSettings.Fail2BanSocketPath
		}
		prometheus.MustRegister(exporter)

		http.Handle("/metrics", promhttp.Handler())
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appSettings.MetricsPort), nil))
	}
}
