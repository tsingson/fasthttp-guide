package webserver

import (
	"net"

	"github.com/valyala/fasthttp/reuseport"
	"go.uber.org/zap"
)

func listen(addr string, log *zap.Logger) (ln net.Listener, err error) {
	ln, err = reuseport.Listen("tcp4", addr)
	if err != nil {
		log.Info("working in windows" + addr)
		// for windows
		ln, err = net.Listen("tcp", addr)
		if err != nil {
			log.Fatal("tcp Connect Error", zap.Error(err))
			return nil, err
		}
	}
	return ln, nil
}
