package f2b

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

const (
	deprecatedNamespace = "fail2ban"
)

var (
	deprecatedMetricUp = prometheus.NewDesc(
		prometheus.BuildFQName(deprecatedNamespace, "", "up"),
		"(Deprecated) Was the last fail2ban query successful.",
		nil, nil,
	)
	deprecatedMetricBannedIpsPerJail = prometheus.NewDesc(
		prometheus.BuildFQName(deprecatedNamespace, "", "banned_ips"),
		"(Deprecated) Number of banned IPs stored in the database (per jail).",
		[]string{"jail"}, nil,
	)
	deprecatedMetricBadIpsPerJail = prometheus.NewDesc(
		prometheus.BuildFQName(deprecatedNamespace, "", "bad_ips"),
		"(Deprecated) Number of bad IPs stored in the database (per jail).",
		[]string{"jail"}, nil,
	)
	deprecatedMetricEnabledJails = prometheus.NewDesc(
		prometheus.BuildFQName(deprecatedNamespace, "", "enabled_jails"),
		"(Deprecated) Enabled jails.",
		[]string{"jail"}, nil,
	)
	deprecatedMetricErrorCount = prometheus.NewDesc(
		prometheus.BuildFQName(deprecatedNamespace, "", "errors"),
		"(Deprecated) Number of errors found since startup.",
		[]string{"type"}, nil,
	)
)

func (c *Collector) collectDeprecatedUpMetric(ch chan<- prometheus.Metric) {
	var upMetricValue float64 = 1
	if c.lastError != nil {
		upMetricValue = 0
	}
	ch <- prometheus.MustNewConstMetric(
		deprecatedMetricUp, prometheus.GaugeValue, upMetricValue,
	)
}

func (c *Collector) collectDeprecatedErrorCountMetric(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		deprecatedMetricErrorCount, prometheus.CounterValue, float64(c.dbErrorCount), "db",
	)
}

func (c *Collector) collectDeprecatedBadIpsPerJailMetrics(ch chan<- prometheus.Metric) {
	jailNameToCountMap, err := c.db.CountBadIpsPerJail()
	c.lastError = err

	if err != nil {
		c.dbErrorCount++
		log.Print(err)
	}

	for jailName, count := range jailNameToCountMap {
		ch <- prometheus.MustNewConstMetric(
			deprecatedMetricBadIpsPerJail, prometheus.GaugeValue, float64(count), jailName,
		)
	}
}

func (c *Collector) collectDeprecatedBannedIpsPerJailMetrics(ch chan<- prometheus.Metric) {
	jailNameToCountMap, err := c.db.CountBannedIpsPerJail()
	c.lastError = err

	if err != nil {
		c.dbErrorCount++
		log.Print(err)
	}

	for jailName, count := range jailNameToCountMap {
		ch <- prometheus.MustNewConstMetric(
			deprecatedMetricBannedIpsPerJail, prometheus.GaugeValue, float64(count), jailName,
		)
	}
}

func (c *Collector) collectDeprecatedEnabledJailMetrics(ch chan<- prometheus.Metric) {
	jailNameToEnabledMap, err := c.db.JailNameToEnabledValue()
	c.lastError = err

	if err != nil {
		c.dbErrorCount++
		log.Print(err)
	}

	for jailName, count := range jailNameToEnabledMap {
		ch <- prometheus.MustNewConstMetric(
			deprecatedMetricEnabledJails, prometheus.GaugeValue, float64(count), jailName,
		)
	}
}
