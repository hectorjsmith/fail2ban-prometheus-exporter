package server

import (
	"log"
	"net/http"

	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/cfg"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/collector/textfile"
)

func StartServer(
	appSettings *cfg.AppSettings,
	textFileCollector *textfile.Collector,
) chan error {
	http.HandleFunc("/", AuthMiddleware(
		rootHtmlHandler,
		appSettings.AuthProvider,
	))
	http.HandleFunc(metricsPath, AuthMiddleware(
		func(w http.ResponseWriter, r *http.Request) {
			metricHandler(w, r, textFileCollector)
		},
		appSettings.AuthProvider,
	))
	log.Printf("metrics available at '%s'", metricsPath)

	svrErr := make(chan error)
	go func() {
		svrErr <- http.ListenAndServe(appSettings.MetricsAddress, nil)
	}()
	log.Print("ready")
	return svrErr
}
