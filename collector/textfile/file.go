package textfile

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "textfile"

var (
	metricReadError = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "error"),
		"Checks for errors while reading text files",
		[]string{"path"}, nil,
	)
)

func (c *Collector) collectFileContents() {
	files, err := os.ReadDir(c.folderPath)
	if err != nil {
		log.Printf("error reading directory '%s': %v", c.folderPath, err)
		return
	}

	for _, file := range files {
		fileName := file.Name()
		if !strings.HasSuffix(strings.ToLower(fileName), ".prom") {
			continue
		}
		c.fileMap[fileName] = &fileData{
			readErrors: 0,
			fileName:   fileName,
		}

		fullPath := filepath.Join(c.folderPath, fileName)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			c.appendErrorForPath(fileName)
			log.Printf("error reading contents of file '%s': %v", fileName, err)
		}

		c.fileMap[fileName].fileContents = content
	}
}

func (c *Collector) collectFileErrors(ch chan<- prometheus.Metric) {
	for _, f := range c.fileMap {
		ch <- prometheus.MustNewConstMetric(
			metricReadError, prometheus.GaugeValue, float64(f.readErrors), f.fileName,
		)
	}
}
