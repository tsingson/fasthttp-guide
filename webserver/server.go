package webserver

import (
	"net"

	"github.com/fasthttp/router"
	"github.com/oklog/run"
	"github.com/tsingson/fastx/utils"
	"github.com/tsingson/zaplogger"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
	"go.uber.org/zap"
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
	addr   string
	Log    *zap.Logger
	ln     net.Listener
	router *router.Router
	debug  bool
}

// NewServer  new fasthttp webServer
func NewServer(cfg WebConfig) *webServer {

	var path, err = utils.GetCurrentExecDir()
	if err != nil {

	}

	logPath := path + "/Log"
	log := zaplogger.NewZapLog(logPath, "vkmsa", true)

	s := &webServer{
		Config: cfg,
		addr:   ":8091",
		Log:    log,
		router: router.New(),
		debug:  true,
	}
	return s
}

// NewServer  new fasthttp webServer
func DefaultServer() *webServer {
	var log = zaplogger.ConsoleDebug()

	var cfg = Default()

	s := &webServer{
		Config: cfg,
		addr:   ServerAddr,
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
	ws.ln, err = listen(ws.addr, ws.Log)
	if err != nil {
		return err
	}
	var lg = zaplogger.InitZapLogger(ws.Log)
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

	// run fasthttp serv
	var g run.Group
	g.Add(func() error {
		return s.Serve(ws.ln)
	}, func(e error) {
		_ = ws.ln.Close()

	})
	return g.Run()
}

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

// design and code by tsingson
