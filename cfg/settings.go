package cfg

type AppSettings struct {
	VersionMode           bool
	MetricsAddress        string
	Fail2BanSocketPath    string
	FileCollectorPath     string
	BasicAuthProvider     *hashedBasicAuth
	ExitOnSocketConnError bool
}
