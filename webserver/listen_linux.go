// +build !windows

package webserver

import (
	"net"
	"runtime"

	"github.com/valyala/fasthttp/reuseport"
	"go.uber.org/zap"
)

func (s *webServer) getListener() (net.Listener, error) {
	if runtime.NumCPU() > 1 {
		ln, err := reuseport.Listen("tcp4", s.Cfg.Addr)
		if err == nil {
			return ln, nil
		}
		s.Log.Warn("Can not use reuseport, using default Listener", zap.Error(err))
	}

	return net.Listen("tcp4", s.Cfg.Addr)
}
