package socket

import (
	"fmt"
	"github.com/kisielk/og-rek"
	"github.com/nlpodyssey/gopickle/types"
	"log"
	"net"
	"strings"
)

type Fail2BanSocket struct {
	socket  net.Conn
	encoder *ogórek.Encoder
}

type JailStats struct {
	FailedCurrent int
	FailedTotal   int
	BannedCurrent int
	BannedTotal   int
}

func ConnectToSocket(path string) (*Fail2BanSocket, error) {
	c, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}
	return &Fail2BanSocket{
		socket:  c,
		encoder: ogórek.NewEncoder(c),
	}, nil
}

func (s *Fail2BanSocket) Close() error {
	return s.socket.Close()
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
	log.Printf("(%s) unexpected response format - cannot parse: %v", pingCommand, response)
	return false
}

func (s *Fail2BanSocket) GetJails() ([]string, error) {
	response, err := s.sendCommand([]string{statusCommand})
	if err != nil {
		return nil, err
	}

	if lvl1, ok := response.(*types.Tuple); ok {
		if lvl2, ok := lvl1.Get(1).(*types.List); ok {
			if lvl3, ok := lvl2.Get(1).(*types.Tuple); ok {
				if lvl4, ok := lvl3.Get(1).(string); ok {
					splitJails := strings.Split(lvl4, ",")
					return trimSpaceForAll(splitJails), nil
				}
			}
		}
	}
	return nil, newBadFormatError(statusCommand, response)
}

func (s *Fail2BanSocket) GetJailStats(jail string) (JailStats, error) {
	response, err := s.sendCommand([]string{statusCommand, jail})
	if err != nil {
		return JailStats{}, err
	}

	stats := JailStats{
		FailedCurrent: -1,
		FailedTotal:   -1,
		BannedCurrent: -1,
		BannedTotal:   -1,
	}

	if lvl1, ok := response.(*types.Tuple); ok {
		if lvl2, ok := lvl1.Get(1).(*types.List); ok {
			if filter, ok := lvl2.Get(0).(*types.Tuple); ok {
				if filterLvl1, ok := filter.Get(1).(*types.List); ok {
					if filterCurrentTuple, ok := filterLvl1.Get(0).(*types.Tuple); ok {
						if filterCurrent, ok := filterCurrentTuple.Get(1).(int); ok {
							stats.FailedCurrent = filterCurrent
						}
					}
					if filterTotalTuple, ok := filterLvl1.Get(1).(*types.Tuple); ok {
						if filterTotal, ok := filterTotalTuple.Get(1).(int); ok {
							stats.FailedTotal = filterTotal
						}
					}
				}
			}
			if actions, ok := lvl2.Get(1).(*types.Tuple); ok {
				if actionsLvl1, ok := actions.Get(1).(*types.List); ok {
					if actionsCurrentTuple, ok := actionsLvl1.Get(0).(*types.Tuple); ok {
						if actionsCurrent, ok := actionsCurrentTuple.Get(1).(int); ok {
							stats.BannedCurrent = actionsCurrent
						}
					}
					if actionsTotalTuple, ok := actionsLvl1.Get(1).(*types.Tuple); ok {
						if actionsTotal, ok := actionsTotalTuple.Get(1).(int); ok {
							stats.BannedTotal = actionsTotal
						}
					}
				}
			}
			return stats, nil
		}
	}
	return stats, newBadFormatError(statusCommand, response)
}

func newBadFormatError(command string, data interface{}) error {
	return fmt.Errorf("(%s) unexpected response format - cannot parse: %v", command, data)
}

func trimSpaceForAll(slice []string) []string {
	for i := range slice {
		slice[i] = strings.TrimSpace(slice[i])
	}
	return slice
}
