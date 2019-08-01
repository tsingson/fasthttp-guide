package webserver

func (ws *webServer) muxRouter() {
	ws.router.GET("/", ws.Recovery(ws.hello()))
	ws.router.GET("/get", ws.Recovery(ws.testGet()))
	ws.router.POST("/post", ws.Recovery(ws.testPost()))
}

// design and code by tsingson
