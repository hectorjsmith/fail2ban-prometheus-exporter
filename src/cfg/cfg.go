package cfg

import (
	"flag"
	"fmt"
	"os"
)

const (
	minServerPort = 1000
	maxServerPort = 65535
)

type AppSettings struct {
	VersionMode    bool
	MetricsPort    int
	Fail2BanDbPath string
}

func Parse() *AppSettings {
	appSettings := &AppSettings{}
	flag.BoolVar(&appSettings.VersionMode, "version", false, "show version info and exit")
	flag.IntVar(&appSettings.MetricsPort, "port", 9191, "port to use for the metrics server")
	flag.StringVar(&appSettings.Fail2BanDbPath, "db", "", "path to the fail2ban sqlite database")

	flag.Parse()
	appSettings.validateFlags()
	return appSettings
}

func (settings *AppSettings) validateFlags() {
	var flagsValid = true
	if !settings.VersionMode {
		if settings.Fail2BanDbPath == "" {
			fmt.Println("missing flag 'db'")
			flagsValid = false
		}
		if settings.MetricsPort < minServerPort || settings.MetricsPort > maxServerPort {
			fmt.Printf("invalid server port, must be within %d and %d (found %d)\n",
				minServerPort, maxServerPort, settings.MetricsPort)
			flagsValid = false
		}
	}
	if !flagsValid {
		flag.Usage()
		os.Exit(1)
	}
}
