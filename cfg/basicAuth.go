package cfg

import "gitlab.com/hectorjsmith/fail2ban-prometheus-exporter/auth"

type hashedBasicAuth struct {
	username string
	password string
}

func newHashedBasicAuth(rawUsername, rawPassword string) *hashedBasicAuth {
	return &hashedBasicAuth{
		username: auth.HashString(rawUsername),
		password: auth.HashString(rawPassword),
	}
}

func (p *hashedBasicAuth) Enabled() bool {
	return len(p.username) > 0 && len(p.password) > 0
}

func (p *hashedBasicAuth) DoesBasicAuthMatch(rawUsername, rawPassword string) bool {
	username := auth.HashString(rawUsername)
	password := auth.HashString(rawPassword)
	return username == p.username && password == p.password
}
