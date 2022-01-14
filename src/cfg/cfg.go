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
	BasicAuthProvider    *hashedBasicAuth
}

func Parse() *AppSettings {
	var rawBasicAuthUsername string
	var rawBasicAuthPassword string

	appSettings := &AppSettings{}
	flag.BoolVar(&appSettings.VersionMode, "version", false, "show version info and exit")
	flag.StringVar(&appSettings.MetricsAddress, "web.listen-address", "0.0.0.0", "address to use for the metrics server")
	flag.IntVar(&appSettings.MetricsPort, "port", 9191, "port to use for the metrics server")
	flag.StringVar(&appSettings.Fail2BanSocketPath, "socket", "", "path to the fail2ban server socket")
	flag.BoolVar(&appSettings.FileCollectorEnabled, "collector.textfile", false, "enable the textfile collector")
	flag.StringVar(&appSettings.FileCollectorPath, "collector.textfile.directory", "", "directory to read text files with metrics from")
	flag.StringVar(&rawBasicAuthUsername, "web.basic-auth.username", "", "username to use to protect endpoints with basic auth")
	flag.StringVar(&rawBasicAuthPassword, "web.basic-auth.password", "", "password to use to protect endpoints with basic auth")

	flag.Parse()
	appSettings.setBasicAuthValues(rawBasicAuthUsername, rawBasicAuthPassword)
	appSettings.validateFlags()
	return appSettings
}

func (settings *AppSettings) setBasicAuthValues(rawUsername, rawPassword string) {
	settings.BasicAuthProvider = newHashedBasicAuth(rawUsername, rawPassword)
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
		if (len(settings.BasicAuthProvider.username) > 0) != (len(settings.BasicAuthProvider.password) > 0) {
			fmt.Printf("to enable basic auth both the username and the password must be provided")
			flagsValid = false
		}
	}
	if !flagsValid {
		flag.Usage()
		os.Exit(1)
	}
}
