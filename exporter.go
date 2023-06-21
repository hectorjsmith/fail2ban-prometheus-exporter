package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/cfg"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/collector/f2b"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/collector/textfile"
	"gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/server"
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

		svrErr := server.StartServer(appSettings, textFileCollector)
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
