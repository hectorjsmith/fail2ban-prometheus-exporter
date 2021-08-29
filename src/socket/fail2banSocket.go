package socket

import (
	"log"
	"net"

	"github.com/kisielk/og-rek"
	"github.com/nlpodyssey/gopickle/types"
)

type Fail2BanSocket struct {
	socket  net.Conn
	encoder *ogórek.Encoder
}

func MustConnectToSocket(path string) *Fail2BanSocket {
	c, err := net.Dial("unix", path)
	if err != nil {
		log.Fatalf("failed to open fail2ban socket: %v", err)
	}
	return &Fail2BanSocket{
		socket:  c,
		encoder: ogórek.NewEncoder(c),
	}
}

func (s *Fail2BanSocket) Ping() bool {
	response, err := s.sendCommand([]string{pingCommand, "100"})
	if err != nil {
		log.Printf("server ping failed: %v", err)
		return false
	}

	if t, ok := response.(*types.Tuple); ok {
		if (*t)[1] == "pong" {
			return true
		}
		log.Printf("unexpected response data: %s", t)
	}
	log.Printf("unexpected response format - cannot parse: %v", response)
	return false
}
