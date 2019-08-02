// +build windows

package webserver

import (
	"net"
)

func (s *webServer) getListener() (net.Listener, error) {
	return net.Listen("tcp4", s.Cfg.Addr)
}
