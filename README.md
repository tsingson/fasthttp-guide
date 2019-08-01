# fasthttp-example


## 0. reference 

1. [GopherCon 2019 - How I write HTTP web services after eight years](https://about.sourcegraph.com/go/gophercon-2019-how-i-write-http-web-services-after-eight-years)


## 1. fasthttp server side

define fasthttp server
```

type webServer struct {
	Config WebConfig
	addr   string
	Log    *zap.Logger
	ln     net.Listener
	router *router.Router
	debug  bool
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

```

router
```
func (ws *webServer) muxRouter() {
	ws.router.GET("/", recovery(ws.hello()))
	ws.router.GET("/get", recovery(ws.testGet()))
	ws.router.POST("/post", recovery(ws.testPost()))
}
```

a POST handler via fasthttp
```
func (ws *webServer) testPost() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		var tid = strconv.FormatInt(int64(ctx.ID()), 10)
		l := ws.Log.Named(tid)
		l.Debug("testPost")

		if ws.debug {
			l.Debug(tid, zap.String("request", ctx.String()))
			ctx.Request.Header.VisitAll(func(key, value []byte) {
				l.Debug(tid, zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
			})
			l.Debug(tid, zap.String("http payload", gotils.B2S(ctx.Request.Body())))
		}

		ctx.SetContentType(ContentRest)
		ctx.SetStatusCode(200)
		ctx.SetBody([]byte(`{"id":2101127497763529765,"plan_start_date":"2019-07-29","plan_end_date":"2019-02-12","actual_start_date":"2019-07-29","actual_end_date":"2019-07-29","is_done":false,"last_updated":"2019-08-01T14:12:17.983236","is_deleted":false,"user_id":2098735545843717147,"title":"00002"}`))
		return
	}
}

```

running fasthttp web-server
```
	var s = webserver.DefaultServer()
...
	s.Run()
	
```


## 2. fasthttp client side

web client visa fasthttp
```

type WebClient struct {
	BaseURI        string
	TransactionID  string
	Authentication bool
	JwtToken       string
	UserAgent      string
	ContentType    string
	Accept         string
	TimeOut        time.Duration
	log            *zap.Logger
	Debug          bool
}

// Default  setup a default fasthttp client
func Default() *WebClient {
	var log = zaplogger.ConsoleDebug()
	return &WebClient{
		Authentication: false,
		TransactionID:  time.Now().String(),
		UserAgent:      "testAgent",
		ContentType:    "application/json; charset=utf-8",
		Accept:         AcceptJson,
		Debug:          true,
		log:            log,
	}
}
```

a GET client
```

// FastGet do GET request via fasthttp
func (w *WebClient) FastGet(requestURI string) (*fasthttp.Response, error) {
	var log = w.log.Named("FastGet")
	t1 := time.Now()
	w.TransactionID = t1.String()
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()
	req.SetRequestURI(requestURI)
	req.Header.SetContentType(w.ContentType)
	req.Header.Add("User-Agent", w.UserAgent)
	req.Header.Add("TransactionID", w.TransactionID)
	req.Header.Add("Accept", w.Accept)
	if w.Authentication && len(w.JwtToken) > 0 {
		req.Header.Set("Authorization", "Bearer "+w.JwtToken)
	}

	// define web client request Method
	req.Header.SetMethod("GET")
	
	
	if w.Debug {
		req.Header.VisitAll(func(key, value []byte) {
			log.Debug(w.TransactionID, zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
		})
		log.Debug(w.TransactionID)
	}

	var timeOut = 3 * time.Second
	if w.TimeOut != 0 {
		timeOut = w.TimeOut
	}
	// DO GET request
	var err = fasthttp.DoTimeout(req, resp, timeOut)
	
	if err != nil {
		log.Error("post request error", zap.Error(err))
		return nil, err
	}
	if w.Debug {
		elapsed := time.Since(t1)
		log.Debug(w.TransactionID, zap.Duration("elapsed", elapsed))
		log.Debug(w.TransactionID, zap.Int("http status code", resp.StatusCode()))
		resp.Header.VisitAll(func(key, value []byte) {
			log.Debug(w.TransactionID, zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
		})
		log.Debug(w.TransactionID, zap.String("http payload", gotils.B2S(resp.Body())))
	}

	// add your logic code here 
	
	var out = fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}
```


call  GET
```
	var w = webclient.Default()
	w.Debug = true

	w.Authentication = true
	w.JwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY2NzIwMDAsInJvbGUiOiJ0ZXJtaW5hbF9hcGsiLCJzdGF0dXMiOiJhY3RpdmUiLCJ1c2VyX2lkIjoiNTBjNjg5MTAtNjEyYi00NjMzLTk2YjktNTA3NzhjNDViNTAwIn0.l1JHnOL85s3ajto0MKs-D6paW1YxpaMuxA0nzI0Xlfk"
	var url = "http://localhost:3001/get"
	var resp, err = w.FastGet(url)

	if err != nil {

	}
	if resp != nil {
		litter.Dump(gotils.B2S(resp.Body()))
	}
	// clean-up
	fasthttp.ReleaseResponse(resp)
```


##  3. clone,  build and run



**build**
```
go install  -gcflags=-trimpath=OPATH -asmflags=-trimpath=OPATH -ldflags "-w -s" ./cmd/... 
```

**run** in two terminal

terminal 1
```
cd $GOBIN
./test-server
```
terminal 2
```
cd $GOBIN
./test-client
```

**output**

client side
```
2019-08-02T06:23:38.671+0800	DEBUG	FastGet	webclient/client.go:126	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"key": "Content-Type", "value": "application/json; charset=utf-8"}
2019-08-02T06:23:38.671+0800	DEBUG	FastGet	webclient/client.go:126	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"key": "User-Agent", "value": "testAgent"}
2019-08-02T06:23:38.672+0800	DEBUG	FastGet	webclient/client.go:126	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"key": "Transactionid", "value": "2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696"}
2019-08-02T06:23:38.672+0800	DEBUG	FastGet	webclient/client.go:126	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"key": "Accept", "value": "application/json"}
2019-08-02T06:23:38.672+0800	DEBUG	FastGet	webclient/client.go:126	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"key": "Authorization", "value": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY2NzIwMDAsInJvbGUiOiJ0ZXJtaW5hbF9hcGsiLCJzdGF0dXMiOiJhY3RpdmUiLCJ1c2VyX2lkIjoiNTBjNjg5MTAtNjEyYi00NjMzLTk2YjktNTA3NzhjNDViNTAwIn0.l1JHnOL85s3ajto0MKs-D6paW1YxpaMuxA0nzI0Xlfk"}
2019-08-02T06:23:38.672+0800	DEBUG	FastGet	webclient/client.go:128	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696
2019-08-02T06:23:38.676+0800	DEBUG	FastGet	webclient/client.go:142	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"elapsed": "4.627531ms"}
2019-08-02T06:23:38.676+0800	DEBUG	FastGet	webclient/client.go:143	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"http status code": 200}
2019-08-02T06:23:38.676+0800	DEBUG	FastGet	webclient/client.go:145	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"key": "Content-Length", "value": "275"}
2019-08-02T06:23:38.676+0800	DEBUG	FastGet	webclient/client.go:145	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"key": "Content-Type", "value": "application/vnd.pgrst.object+json; charset=utf-8"}
2019-08-02T06:23:38.676+0800	DEBUG	FastGet	webclient/client.go:145	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"key": "Server", "value": "EPG-xcache-service"}
2019-08-02T06:23:38.676+0800	DEBUG	FastGet	webclient/client.go:145	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"key": "Date", "value": "Thu, 01 Aug 2019 22:23:37 GMT"}
2019-08-02T06:23:38.676+0800	DEBUG	FastGet	webclient/client.go:147	2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696	{"http payload": "{\"id\":2101127497763529765,\"plan_start_date\":\"2019-07-29\",\"plan_end_date\":\"2019-02-12\",\"actual_start_date\":\"2019-07-29\",\"actual_end_date\":\"2019-07-29\",\"is_done\":false,\"last_updated\":\"2019-08-01T14:12:17.983236\",\"is_deleted\":false,\"user_id\":2098735545843717147,\"title\":\"00002\"}"}
"{\"id\":2101127497763529765,\"plan_start_date\":\"2019-07-29\",\"plan_end_date\":\"2019-02-12\",\"actual_start_date\":\"2019-07-29\",\"actual_end_date\":\"2019-07-29\",\"is_done\":false,\"last_updated\":\"2019-08-01T14:12:17.983236\",\"is_deleted\":false,\"user_id\":2098735545843717147,\"title\":\"00002\"}"
2019-08-02T06:23:38.677+0800	DEBUG	FastPostByte	webclient/client.go:71	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961	{"key": "Content-Type", "value": "application/json; charset=utf-8"}
2019-08-02T06:23:38.677+0800	DEBUG	FastPostByte	webclient/client.go:71	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961	{"key": "User-Agent", "value": "testAgent"}
2019-08-02T06:23:38.677+0800	DEBUG	FastPostByte	webclient/client.go:71	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961	{"key": "Transactionid", "value": "2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961"}
2019-08-02T06:23:38.677+0800	DEBUG	FastPostByte	webclient/client.go:71	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961	{"key": "Accept", "value": "application/vnd.pgrst.object+json"}
2019-08-02T06:23:38.677+0800	DEBUG	FastPostByte	webclient/client.go:73	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961
2019-08-02T06:23:38.678+0800	DEBUG	FastPostByte	webclient/client.go:87	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961	{"elapsed": "1.19083ms"}
2019-08-02T06:23:38.678+0800	DEBUG	FastPostByte	webclient/client.go:88	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961	{"http status code": 200}
2019-08-02T06:23:38.678+0800	DEBUG	FastPostByte	webclient/client.go:90	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961	{"key": "Content-Length", "value": "275"}
2019-08-02T06:23:38.678+0800	DEBUG	FastPostByte	webclient/client.go:90	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961	{"key": "Content-Type", "value": "application/vnd.pgrst.object+json; charset=utf-8"}
2019-08-02T06:23:38.678+0800	DEBUG	FastPostByte	webclient/client.go:90	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961	{"key": "Server", "value": "EPG-xcache-service"}
2019-08-02T06:23:38.678+0800	DEBUG	FastPostByte	webclient/client.go:90	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961	{"key": "Date", "value": "Thu, 01 Aug 2019 22:23:37 GMT"}
2019-08-02T06:23:38.678+0800	DEBUG	FastPostByte	webclient/client.go:92	2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961	{"http payload": "{\"id\":2101127497763529765,\"plan_start_date\":\"2019-07-29\",\"plan_end_date\":\"2019-02-12\",\"actual_start_date\":\"2019-07-29\",\"actual_end_date\":\"2019-07-29\",\"is_done\":false,\"last_updated\":\"2019-08-01T14:12:17.983236\",\"is_deleted\":false,\"user_id\":2098735545843717147,\"title\":\"00002\"}"}
"{\"id\":2101127497763529765,\"plan_start_date\":\"2019-07-29\",\"plan_end_date\":\"2019-02-12\",\"actual_start_date\":\"2019-07-29\",\"actual_end_date\":\"2019-07-29\",\"is_done\":false,\"last_updated\":\"2019-08-01T14:12:17.983236\",\"is_deleted\":false,\"user_id\":2098735545843717147,\"title\":\"00002\"}"
```

**server side**
```
----- fasthttp server starting -----
2019-08-02T06:23:38.675+0800	DEBUG	4294967297	webserver/testHandler.go:15	testGet
2019-08-02T06:23:38.675+0800	DEBUG	4294967297	webserver/testHandler.go:19	4294967297	{"request": "#0000000100000001 - 127.0.0.1:3001<->127.0.0.1:64674 - GET http://localhost:3001/get"}
2019-08-02T06:23:38.675+0800	DEBUG	4294967297	webserver/testHandler.go:22	4294967297	{"key": "Host", "value": "localhost:3001"}
2019-08-02T06:23:38.675+0800	DEBUG	4294967297	webserver/testHandler.go:22	4294967297	{"key": "Content-Length", "value": "0"}
2019-08-02T06:23:38.675+0800	DEBUG	4294967297	webserver/testHandler.go:22	4294967297	{"key": "Content-Type", "value": "application/json; charset=utf-8"}
2019-08-02T06:23:38.675+0800	DEBUG	4294967297	webserver/testHandler.go:22	4294967297	{"key": "User-Agent", "value": "testAgent"}
2019-08-02T06:23:38.675+0800	DEBUG	4294967297	webserver/testHandler.go:22	4294967297	{"key": "Transactionid", "value": "2019-08-02 06:23:38.671678 +0800 CST m=+0.004252696"}
2019-08-02T06:23:38.675+0800	DEBUG	4294967297	webserver/testHandler.go:22	4294967297	{"key": "Accept", "value": "application/json"}
2019-08-02T06:23:38.675+0800	DEBUG	4294967297	webserver/testHandler.go:22	4294967297	{"key": "Authorization", "value": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY2NzIwMDAsInJvbGUiOiJ0ZXJtaW5hbF9hcGsiLCJzdGF0dXMiOiJhY3RpdmUiLCJ1c2VyX2lkIjoiNTBjNjg5MTAtNjEyYi00NjMzLTk2YjktNTA3NzhjNDViNTAwIn0.l1JHnOL85s3ajto0MKs-D6paW1YxpaMuxA0nzI0Xlfk"}
2019-08-02T06:23:38.676+0800	DEBUG	4294967297	webserver/testHandler.go:25	4294967297	{"http payload": ""}
2019-08-02T06:23:38.677+0800	DEBUG	4294967298	webserver/testHandler.go:42	testPost
2019-08-02T06:23:38.677+0800	DEBUG	4294967298	webserver/testHandler.go:46	4294967298	{"request": "#0000000100000002 - 127.0.0.1:3001<->127.0.0.1:64674 - POST http://localhost:3001/post"}
2019-08-02T06:23:38.677+0800	DEBUG	4294967298	webserver/testHandler.go:50	4294967298	{"key": "Host", "value": "localhost:3001"}
2019-08-02T06:23:38.677+0800	DEBUG	4294967298	webserver/testHandler.go:50	4294967298	{"key": "Content-Length", "value": "183"}
2019-08-02T06:23:38.677+0800	DEBUG	4294967298	webserver/testHandler.go:50	4294967298	{"key": "Content-Type", "value": "application/json; charset=utf-8"}
2019-08-02T06:23:38.677+0800	DEBUG	4294967298	webserver/testHandler.go:50	4294967298	{"key": "User-Agent", "value": "testAgent"}
2019-08-02T06:23:38.677+0800	DEBUG	4294967298	webserver/testHandler.go:50	4294967298	{"key": "Transactionid", "value": "2019-08-02 06:23:38.677108 +0800 CST m=+0.009682961"}
2019-08-02T06:23:38.677+0800	DEBUG	4294967298	webserver/testHandler.go:50	4294967298	{"key": "Accept", "value": "application/vnd.pgrst.object+json"}
2019-08-02T06:23:38.677+0800	DEBUG	4294967298	webserver/testHandler.go:53	4294967298	{"http payload": "{\"actual_start_date\":\"2019-07-29\",\"actual_end_date\":\"2019-07-29\",\"plan_start_date\":\"2019-07-29\",\"plan_end_date\":\"2019-02-12\",\"title\":\"养殖计划00002\",\"user_id\":2098735545843717147}"}
```


