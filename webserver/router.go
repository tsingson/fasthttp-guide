package webserver

func (ws *webServer) muxRouter() {
	ws.router.GET("/", recovery(ws.hello()))
	ws.router.GET("/get", recovery(ws.testGet()))
	ws.router.POST("/post", recovery(ws.testPost()))
}

// design and code by tsingson
