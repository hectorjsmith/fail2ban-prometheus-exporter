package export

import (
	"fail2ban-prometheus-exporter/cfg"
	fail2banDb "fail2ban-prometheus-exporter/db"
	"log"
)

func NewExporter(appSettings *cfg.AppSettings, exporterVersion string) *Exporter {
	exporter := &Exporter{
		exporterVersion:            exporterVersion,
		lastError:                  nil,
		dbErrorCount:               0,
		socketConnectionErrorCount: 0,
		socketRequestErrorCount:    0,
	}
	if appSettings.Fail2BanDbPath != "" {
		log.Print("database-based metrics have been deprecated and will be removed in a future release")
		exporter.db = fail2banDb.MustConnectToDb(appSettings.Fail2BanDbPath)
	}
	if appSettings.Fail2BanSocketPath != "" {
		exporter.socketPath = appSettings.Fail2BanSocketPath
	}
	return exporter
}
