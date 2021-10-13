package f2b

import (
	"fail2ban-prometheus-exporter/cfg"
	fail2banDb "fail2ban-prometheus-exporter/db"
	"fail2ban-prometheus-exporter/socket"
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

type Collector struct {
	db                         *fail2banDb.Fail2BanDB
	socketPath                 string
	exporterVersion            string
	lastError                  error
	dbErrorCount               int
	socketConnectionErrorCount int
	socketRequestErrorCount    int
}

func NewExporter(appSettings *cfg.AppSettings, exporterVersion string) *Collector {
	colector := &Collector{
		exporterVersion:            exporterVersion,
		lastError:                  nil,
		dbErrorCount:               0,
		socketConnectionErrorCount: 0,
		socketRequestErrorCount:    0,
	}
	if appSettings.Fail2BanDbPath != "" {
		log.Print("database-based metrics have been deprecated and will be removed in a future release")
		colector.db = fail2banDb.MustConnectToDb(appSettings.Fail2BanDbPath)
	}
	if appSettings.Fail2BanSocketPath != "" {
		colector.socketPath = appSettings.Fail2BanSocketPath
	}
	return colector
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	if c.db != nil {
		ch <- deprecatedMetricUp
		ch <- deprecatedMetricBadIpsPerJail
		ch <- deprecatedMetricBannedIpsPerJail
		ch <- deprecatedMetricEnabledJails
		ch <- deprecatedMetricErrorCount
	}
	if c.socketPath != "" {
		ch <- metricServerUp
		ch <- metricJailCount
		ch <- metricJailFailedCurrent
		ch <- metricJailFailedTotal
		ch <- metricJailBannedCurrent
		ch <- metricJailBannedTotal
	}
	ch <- metricErrorCount
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	if c.db != nil {
		c.collectDeprecatedBadIpsPerJailMetrics(ch)
		c.collectDeprecatedBannedIpsPerJailMetrics(ch)
		c.collectDeprecatedEnabledJailMetrics(ch)
		c.collectDeprecatedUpMetric(ch)
		c.collectDeprecatedErrorCountMetric(ch)
	}
	if c.socketPath != "" {
		s, err := socket.ConnectToSocket(c.socketPath)
		if err != nil {
			log.Printf("error opening socket: %v", err)
			c.socketConnectionErrorCount++
		} else {
			defer s.Close()
		}
		c.collectServerUpMetric(ch, s)
		if err == nil && s != nil {
			c.collectJailMetrics(ch, s)
		}
		c.collectVersionMetric(ch, s)
	} else {
		c.collectVersionMetric(ch, nil)
	}
	c.collectErrorCountMetric(ch)
}
