package webserver

import (
	"net"
	"sync"

	"github.com/fasthttp/router"
	"github.com/oklog/run"
	"github.com/tsingson/fastx/utils"
	"github.com/tsingson/zaplogger"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const (
	AcceptJson  = "application/json"
	AcceptRest  = "application/vnd.pgrst.object+json"
	ContentText = "text/plain; charset=utf8"
	ContentRest = "application/vnd.pgrst.object+json; charset=utf-8"
	ContentJson = "application/json; charset=utf-8"
)

var once sync.Once

type webServer struct {
	Cfg    WebConfig
	addr   string
	Log    *zap.Logger
	ln     []net.Listener
	router *router.Router
	debug  bool
}

// NewServer  new fasthttp webServer
func NewServer(cfg WebConfig) (s *webServer) {
	once.Do(func() {
		var path, _ = utils.GetCurrentExecDir()

		logPath := path + "/Log"
		var log = zaplogger.NewZapLog(logPath, "vkmsa", true)

		s = &webServer{
			Cfg:    cfg,
			addr:   ":8091",
			Log:    log,
			router: router.New(),
			debug:  true,
		}
	})
	return s
}

// NewServer  new fasthttp webServer
func DefaultServer() (s *webServer) {
	once.Do(func() {
		var log = zaplogger.ConsoleDebug()
		var cfg = Default()
		s = &webServer{
			Cfg:    cfg,
			addr:   ServerAddr,
			Log:    log,
			router: router.New(),
			debug:  true,
		}
	})
	return s
}

func (ws *webServer) Close() {
	for _, v := range ws.ln {
		_ = v.Close()
	}
}

func (ws *webServer) Run() (err error) {
	ws.muxRouter()
	// reuse port

	// for i:=0; i< runtime.NumCPU(); i++ {
	var ln net.Listener
	// var err error
	ln, err = ws.getListener()
	if err != nil {
		return err
	}
	ws.ln = append(ws.ln, ln)
	// }

	var lg = zaplogger.InitZapLogger(ws.Log)
	s := &fasthttp.Server{
		Handler:            ws.router.Handler,
		Name:               ws.Cfg.Name,
		ReadBufferSize:     ws.Cfg.ReadBufferSize,
		MaxConnsPerIP:      ws.Cfg.MaxConnsPerIP,
		MaxRequestsPerConn: ws.Cfg.MaxRequestsPerConn,
		MaxRequestBodySize: ws.Cfg.MaxRequestBodySize, //  100 << 20, // 100MB // 1024 * 4, // MaxRequestBodySize:
		Concurrency:        ws.Cfg.Concurrency,
		Logger:             lg,
	}

	// run fasthttp serv
	var g run.Group
	for _, v := range ws.ln {
		g.Add(func() error {
			return s.Serve(v)
		}, func(e error) {
			_ = v.Close()
		})
	}
	return g.Run()
}

// design and code by tsingson
