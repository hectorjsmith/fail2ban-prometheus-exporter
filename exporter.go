package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/auth"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/cfg"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/collector/f2b"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/collector/textfile"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	metricsPath = "/metrics"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func printAppVersion() {
	fmt.Println(version)
	fmt.Printf("    build date:  %s\r\n    commit hash: %s\r\n    built by:    %s\r\n", date, commit, builtBy)
}

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

func main() {
	appSettings := cfg.Parse()
	if appSettings.VersionMode {
		printAppVersion()
	} else {
		handleGracefulShutdown()
		log.Printf("fail2ban exporter version %s", version)
		log.Printf("starting server at %s", appSettings.MetricsAddress)

		f2bCollector := f2b.NewExporter(appSettings, version)
		prometheus.MustRegister(f2bCollector)

		textFileCollector := textfile.NewCollector(appSettings)
		prometheus.MustRegister(textFileCollector)

		http.HandleFunc("/", auth.BasicAuthMiddleware(rootHtmlHandler, appSettings.BasicAuthProvider))
		http.HandleFunc(metricsPath, auth.BasicAuthMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				metricHandler(w, r, textFileCollector)
			},
			appSettings.BasicAuthProvider,
		))
		log.Printf("metrics available at '%s'", metricsPath)
		if appSettings.BasicAuthProvider.Enabled() {
			log.Printf("basic auth enabled")
		}

		svrErr := make(chan error)
		go func() {
			svrErr <- http.ListenAndServe(appSettings.MetricsAddress, nil)
		}()
		log.Print("ready")

		err := <-svrErr
		log.Fatal(err)
	}
}

func handleGracefulShutdown() {
	var signals = make(chan os.Signal)

	signal.Notify(signals, syscall.SIGTERM)
	signal.Notify(signals, syscall.SIGINT)

	go func() {
		sig := <-signals
		log.Printf("caught signal: %+v", sig)
		os.Exit(0)
	}()
}
