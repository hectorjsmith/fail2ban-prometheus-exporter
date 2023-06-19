package socket

import (
	"fmt"
	"github.com/kisielk/og-rek"
	"github.com/nlpodyssey/gopickle/types"
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

func (s *Fail2BanSocket) Ping() (bool, error) {
	response, err := s.sendCommand([]string{pingCommand, "100"})
	if err != nil {
		return false, newConnectionError(pingCommand, err)
	}

	if t, ok := response.(*types.Tuple); ok {
		if (*t)[1] == "pong" {
			return true, nil
		}
		return false, fmt.Errorf("unexpected response data (expecting 'pong'): %s", (*t)[1])
	}
	return false, newBadFormatError(pingCommand, response)
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

func (s *Fail2BanSocket) GetJailBanTime(jail string) (int, error) {
	command := fmt.Sprintf(banTimeCommandFmt, jail)
	return s.sendSimpleIntCommand(command)
}

func (s *Fail2BanSocket) GetJailFindTime(jail string) (int, error) {
	command := fmt.Sprintf(findTimeCommandFmt, jail)
	return s.sendSimpleIntCommand(command)
}

func (s *Fail2BanSocket) GetJailMaxRetries(jail string) (int, error) {
	command := fmt.Sprintf(maxRetriesCommandFmt, jail)
	return s.sendSimpleIntCommand(command)
}

func (s *Fail2BanSocket) GetServerVersion() (string, error) {
	response, err := s.sendCommand([]string{versionCommand})
	if err != nil {
		return "", err
	}

	if lvl1, ok := response.(*types.Tuple); ok {
		if versionStr, ok := lvl1.Get(1).(string); ok {
			return versionStr, nil
		}
	}
	return "", newBadFormatError(versionCommand, response)
}

// sendSimpleIntCommand sends a command to the fail2ban socket and parses the response to extract an int.
// This command assumes that the response data is in the format of `(d, d)` where `d` is a number.
func (s *Fail2BanSocket) sendSimpleIntCommand(command string) (int, error) {
	response, err := s.sendCommand(strings.Split(command, " "))
	if err != nil {
		return -1, err
	}

	if lvl1, ok := response.(*types.Tuple); ok {
		if banTime, ok := lvl1.Get(1).(int); ok {
			return banTime, nil
		}
	}
	return -1, newBadFormatError(command, response)
}

func newBadFormatError(command string, data interface{}) error {
	return fmt.Errorf("(%s) unexpected response format - cannot parse: %v", command, data)
}

func newConnectionError(command string, err error) error {
	return fmt.Errorf("(%s) failed to send command through socket: %v", command, err)
}

func trimSpaceForAll(slice []string) []string {
	for i := range slice {
		slice[i] = strings.TrimSpace(slice[i])
	}
	return slice
}
