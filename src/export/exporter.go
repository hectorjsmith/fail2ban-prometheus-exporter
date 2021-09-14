package export

import (
	fail2banDb "fail2ban-prometheus-exporter/db"
	"fail2ban-prometheus-exporter/socket"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type Exporter struct {
	db                         *fail2banDb.Fail2BanDB
	socketPath                 string
	exporterVersion            string
	lastError                  error
	dbErrorCount               int
	socketConnectionErrorCount int
	socketRequestErrorCount    int
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	if e.db != nil {
		ch <- deprecatedMetricUp
		ch <- deprecatedMetricBadIpsPerJail
		ch <- deprecatedMetricBannedIpsPerJail
		ch <- deprecatedMetricEnabledJails
		ch <- deprecatedMetricErrorCount
	}
	if e.socketPath != "" {
		ch <- metricServerUp
		ch <- metricJailCount
		ch <- metricJailFailedCurrent
		ch <- metricJailFailedTotal
		ch <- metricJailBannedCurrent
		ch <- metricJailBannedTotal
	}
	ch <- metricErrorCount
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	if e.db != nil {
		e.collectDeprecatedBadIpsPerJailMetrics(ch)
		e.collectDeprecatedBannedIpsPerJailMetrics(ch)
		e.collectDeprecatedEnabledJailMetrics(ch)
		e.collectDeprecatedUpMetric(ch)
		e.collectDeprecatedErrorCountMetric(ch)
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
		e.collectVersionMetric(ch, s)
	} else {
		e.collectVersionMetric(ch, nil)
	}
	e.collectErrorCountMetric(ch)
}
