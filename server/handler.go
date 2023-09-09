package server

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/collector/f2b"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/collector/textfile"
)

const (
	metricsPath = "/metrics"
)

func rootHtmlHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(
		`<html>
			<head><title>Fail2Ban Exporter</title></head>
			<body>
			<h1>Fail2Ban Exporter</h1>
			<p><a href="` + metricsPath + `">Metrics</a></p>
			</body>
		</html>`))
	if err != nil {
		log.Printf("error handling root url: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func metricHandler(w http.ResponseWriter, r *http.Request, collector *textfile.Collector) {
	promhttp.Handler().ServeHTTP(w, r)
	collector.WriteTextFileMetrics(w, r)
}

func healthHandler(w http.ResponseWriter, r *http.Request, collector *f2b.Collector) {
	if collector.IsHealthy() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{\"healthy\":true}"))
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"healthy\":false}"))
	}
}
