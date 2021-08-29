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
	VersionMode        bool
	MetricsPort        int
	Fail2BanDbPath     string
	Fail2BanSocketPath string
}

func Parse() *AppSettings {
	appSettings := &AppSettings{}
	flag.BoolVar(&appSettings.VersionMode, "version", false, "show version info and exit")
	flag.IntVar(&appSettings.MetricsPort, "port", 9191, "port to use for the metrics server")
	flag.StringVar(&appSettings.Fail2BanDbPath, "db", "", "path to the fail2ban sqlite database")
	flag.StringVar(&appSettings.Fail2BanSocketPath, "socket", "", "path to the fail2ban server socket")

	flag.Parse()
	appSettings.validateFlags()
	return appSettings
}

func (settings *AppSettings) validateFlags() {
	var flagsValid = true
	if !settings.VersionMode {
		if settings.Fail2BanDbPath == "" && settings.Fail2BanSocketPath == "" {
			fmt.Println("at least one of the following flags must be provided: 'db', 'socket'")
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
