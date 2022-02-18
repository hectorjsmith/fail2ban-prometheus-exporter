package cfg

import (
	"flag"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

const (
	minServerPort               = 1000
	maxServerPort               = 65535
	socketEnvName               = "F2B_COLLECTOR_SOCKET"
	fileCollectorEnabledEnvName = "F2B_COLLECTOR_TEXT"
	fileCollectorPathEnvName    = "F2B_COLLECTOR_TEXT_PATH"
	portEnvName                 = "F2B_WEB_PORT"
	addressEnvName              = "F2B_WEB_LISTEN_ADDRESS"
	basicAuthUserEnvName        = "F2B_WEB_BASICAUTH_USER"
	basicAuthPassEnvName        = "F2B_WEB_BASICAUTH_PASS"
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

func init() {
	kingpin.HelpFlag.Short('h')
}

func Parse() *AppSettings {
	settings := &AppSettings{}
	readParamsFromCli(settings)
	settings.validateFlags()
	return settings
}

func readParamsFromCli(settings *AppSettings) {
	versionMode := kingpin.
		Flag("version", "show version info and exit").
		Default("false").
		Bool()
	socketPath := kingpin.
		Flag("socket", "path to the fail2ban server socket").
		Default("/var/run/fail2ban/fail2ban.sock").
		Envar(socketEnvName).
		String()
	fileCollectorEnabled := kingpin.
		Flag("collector.textfile", "enable the textfile collector").
		Default("false").
		Envar(fileCollectorEnabledEnvName).
		Bool()
	fileCollectorPath := kingpin.
		Flag("collector.textfile.directory", "directory to read text files with metrics from").
		Default("").
		Envar(fileCollectorPathEnvName).
		String()
	port := kingpin.
		Flag("port", "port to use for the metrics server").
		Default("9191").
		Envar(portEnvName).
		Int()
	address := kingpin.
		Flag("web.listen-address", "address to use for the metrics server").
		Default("0.0.0.0").
		Envar(addressEnvName).
		String()
	rawBasicAuthUsername := kingpin.
		Flag("web.basic-auth.username", "username to use to protect endpoints with basic auth").
		Default("").
		Envar(basicAuthUserEnvName).
		String()
	rawBasicAuthPassword := kingpin.
		Flag("web.basic-auth.password", "password to use to protect endpoints with basic auth").
		Default("").
		Envar(basicAuthPassEnvName).
		String()

	kingpin.Parse()

	settings.VersionMode = *versionMode
	settings.MetricsPort = *port
	settings.MetricsAddress = *address
	settings.Fail2BanSocketPath = *socketPath
	settings.FileCollectorEnabled = *fileCollectorEnabled
	settings.FileCollectorPath = *fileCollectorPath
	settings.setBasicAuthValues(*rawBasicAuthUsername, *rawBasicAuthPassword)
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
