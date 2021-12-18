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
	VersionMode          bool
	MetricsAddress       string
	MetricsPort          int
	Fail2BanSocketPath   string
	FileCollectorPath    string
	FileCollectorEnabled bool
}

func Parse() *AppSettings {
	appSettings := &AppSettings{}
	flag.BoolVar(&appSettings.VersionMode, "version", false, "show version info and exit")
	flag.StringVar(&appSettings.MetricsAddress, "web.listen-address", "0.0.0.0", "address to use for the metrics server")
	flag.IntVar(&appSettings.MetricsPort, "port", 9191, "port to use for the metrics server")
	flag.StringVar(&appSettings.Fail2BanSocketPath, "socket", "", "path to the fail2ban server socket")
	flag.BoolVar(&appSettings.FileCollectorEnabled, "collector.textfile", false, "enable the textfile collector")
	flag.StringVar(&appSettings.FileCollectorPath, "collector.textfile.directory", "", "directory to read text files with metrics from")

	// deprecated: to be removed in next version
	_ = flag.String("db", "", "path to the fail2ban sqlite database (removed)")

	flag.Parse()
	appSettings.validateFlags()
	return appSettings
}

func (settings *AppSettings) validateFlags() {
	var flagsValid = true
	if !settings.VersionMode {
		if settings.Fail2BanSocketPath == "" {
			fmt.Println("fail2ban socket path must not be blank")
			flagsValid = false
		}
		if settings.MetricsPort < minServerPort || settings.MetricsPort > maxServerPort {
			fmt.Printf("invalid server port, must be within %d and %d (found %d)\n",
				minServerPort, maxServerPort, settings.MetricsPort)
			flagsValid = false
		}
		if settings.FileCollectorEnabled && settings.FileCollectorPath == "" {
			fmt.Printf("file collector directory path must not be empty if collector enabled\n")
			flagsValid = false
		}
	}
	if !flagsValid {
		flag.Usage()
		os.Exit(1)
	}
}
