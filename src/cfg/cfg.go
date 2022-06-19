package cfg

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

const (
	socketEnvName                = "F2B_COLLECTOR_SOCKET"
	fileCollectorPathEnvName     = "F2B_COLLECTOR_TEXT_PATH"
	addressEnvName               = "F2B_WEB_LISTEN_ADDRESS"
	basicAuthUserEnvName         = "F2B_WEB_BASICAUTH_USER"
	basicAuthPassEnvName         = "F2B_WEB_BASICAUTH_PASS"
	exitOnSocketConnErrorEnvName = "F2B_EXIT_ON_SOCKET_CONN_ERROR"
)

type AppSettings struct {
	VersionMode           bool
	MetricsAddress        string
	Fail2BanSocketPath    string
	FileCollectorPath     string
	BasicAuthProvider     *hashedBasicAuth
	ExitOnSocketConnError bool
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
		Short('v').
		Default("false").
		Bool()
	socketPath := kingpin.
		Flag("collector.f2b.socket", "path to the fail2ban server socket").
		Default("/var/run/fail2ban/fail2ban.sock").
		Envar(socketEnvName).
		String()
	fileCollectorPath := kingpin.
		Flag("collector.textfile.directory", "directory to read text files with metrics from").
		Default("").
		Envar(fileCollectorPathEnvName).
		String()
	address := kingpin.
		Flag("web.listen-address", "address to use for the metrics server").
		Default(":9191").
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
	rawExitOnSocketConnError := kingpin.
		Flag("collector.f2b.exit-on-socket-connection-error", "when set to true the exporter will immediately exit on a fail2ban socket connection error").
		Default("false").
		Envar(exitOnSocketConnErrorEnvName).
		Bool()

	kingpin.Parse()

	settings.VersionMode = *versionMode
	settings.MetricsAddress = *address
	settings.Fail2BanSocketPath = *socketPath
	settings.FileCollectorPath = *fileCollectorPath
	settings.setBasicAuthValues(*rawBasicAuthUsername, *rawBasicAuthPassword)
	settings.ExitOnSocketConnError = *rawExitOnSocketConnError
}

func (settings *AppSettings) setBasicAuthValues(rawUsername, rawPassword string) {
	settings.BasicAuthProvider = newHashedBasicAuth(rawUsername, rawPassword)
}

func (settings *AppSettings) validateFlags() {
	var flagsValid = true
	if !settings.VersionMode {
		if settings.Fail2BanSocketPath == "" {
			fmt.Println("error: fail2ban socket path must not be blank")
			flagsValid = false
		}
		if settings.MetricsAddress == "" {
			fmt.Println("error: invalid server address, must not be blank")
			flagsValid = false
		}
		if (len(settings.BasicAuthProvider.username) > 0) != (len(settings.BasicAuthProvider.password) > 0) {
			fmt.Println("error: to enable basic auth both the username and the password must be provided")
			flagsValid = false
		}
	}
	if !flagsValid {
		kingpin.Usage()
		os.Exit(1)
	}
}
