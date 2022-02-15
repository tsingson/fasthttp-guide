package webserver

import (
	"net"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp/reuseport"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"

	"github.com/tsingson/fasthttp-guide/logger"
)

const (
	AcceptJson  = "application/json"
	AcceptRest  = "application/vnd.pgrst.object+json"
	ContentText = "text/plain; charset=utf8"
	ContentRest = "application/vnd.pgrst.object+json; charset=utf-8"
	ContentJson = "application/json; charset=utf-8"
)

type webServer struct {
	Config WebConfig
	Addr   string
	Log    *zap.Logger
	ln     net.Listener
	router *router.Router
	debug  bool
}

// NewServer  new fasthttp webServer
func NewServer(cfg WebConfig) *webServer {
	log := logger.Console()

	s := &webServer{
		Config: cfg,
		Addr:   ServerAddr,
		Log:    log,
		router: router.New(),
		debug:  true,
	}
	return s
}

func (ws *webServer) Close() {
	_ = ws.ln.Close()
}

func (ws *webServer) Run() (err error) {
	ws.muxRouter()
	// reuse port

	ws.ln, err = reuseport.Listen("tcp4", ws.Addr)
	if err != nil {
		return err
	}
	lg := logger.InitZapLogger(ws.Log)
	s := &fasthttp.Server{
		Handler:            ws.router.Handler,
		Name:               ws.Config.Name,
		ReadBufferSize:     ws.Config.ReadBufferSize,
		MaxConnsPerIP:      ws.Config.MaxConnsPerIP,
		MaxRequestsPerConn: ws.Config.MaxRequestsPerConn,
		MaxRequestBodySize: ws.Config.MaxRequestBodySize, //  100 << 20, // 100MB // 1024 * 4, // MaxRequestBodySize:
		Concurrency:        ws.Config.Concurrency,
		Logger:             lg,
	}

	return s.Serve(ws.ln)
}

// design and code by tsingson
