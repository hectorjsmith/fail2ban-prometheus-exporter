package export

import (
	"fail2ban-prometheus-exporter/socket"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

const (
	namespace = "f2b"
)

var (
	metricErrorCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "errors"),
		"Number of errors found since startup",
		[]string{"type"}, nil,
	)
	metricServerUp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Check if the fail2ban server is up",
		nil, nil,
	)
	metricJailCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jail_count"),
		"Number of defined jails",
		nil, nil,
	)
	metricJailFailedCurrent = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jail_failed_current"),
		"Number of current failures on this jail's filter",
		[]string{"jail"}, nil,
	)
	metricJailFailedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jail_failed_total"),
		"Number of total failures on this jail's filter",
		[]string{"jail"}, nil,
	)
	metricJailBannedCurrent = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jail_banned_current"),
		"Number of IPs currently banned in this jail",
		[]string{"jail"}, nil,
	)
	metricJailBannedTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "jail_banned_total"),
		"Total number of IPs banned by this jail (includes expired bans)",
		[]string{"jail"}, nil,
	)
	metricVersionInfo = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "version"),
		"Version of the exporter and fail2ban server",
		[]string{"exporter", "fail2ban"}, nil,
	)
)

func (e *Exporter) collectErrorCountMetric(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		metricErrorCount, prometheus.CounterValue, float64(e.dbErrorCount), "db",
	)
	ch <- prometheus.MustNewConstMetric(
		metricErrorCount, prometheus.CounterValue, float64(e.socketConnectionErrorCount), "socket_conn",
	)
	ch <- prometheus.MustNewConstMetric(
		metricErrorCount, prometheus.CounterValue, float64(e.socketRequestErrorCount), "socket_req",
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

func (e *Exporter) collectVersionMetric(ch chan<- prometheus.Metric, s *socket.Fail2BanSocket) {
	var err error
	var fail2banVersion = ""
	if s != nil {
		fail2banVersion, err = s.GetServerVersion()
		if err != nil {
			e.socketRequestErrorCount++
			log.Printf("failed to get fail2ban server version: %v", err)
		}
	}

	ch <- prometheus.MustNewConstMetric(
		metricVersionInfo, prometheus.GaugeValue, float64(1), e.exporterVersion, fail2banVersion,
	)
}
