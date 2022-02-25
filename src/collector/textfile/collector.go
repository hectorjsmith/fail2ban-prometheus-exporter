package textfile

import (
	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/cfg"
	"log"
)

type Collector struct {
	enabled    bool
	folderPath string
	fileMap    map[string]*fileData
}

type fileData struct {
	readErrors   int
	fileName     string
	fileContents []byte
}

func NewCollector(appSettings *cfg.AppSettings) *Collector {
	collector := &Collector{
		enabled:    appSettings.FileCollectorPath != "",
		folderPath: appSettings.FileCollectorPath,
		fileMap:    make(map[string]*fileData),
	}
	if collector.enabled {
		log.Printf("reading textfile metrics from: %s", collector.folderPath)
	}
	return collector
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	if c.enabled {
		ch <- metricReadError
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	if c.enabled {
		c.collectFileContents()
		c.collectFileErrors(ch)
	}
}

func (c *Collector) appendErrorForPath(path string) {
	c.fileMap[path].readErrors++
}
