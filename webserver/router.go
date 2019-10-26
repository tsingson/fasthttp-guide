package webserver

func (ws *webServer) muxRouter() {
	ws.router.GET("/", ws.Recovery(ws.helloWorldGetHandler()))
	ws.router.GET("/get", ws.Recovery(ws.simpleGetHandler()))
	ws.router.POST("/post", ws.Recovery(ws.simplePostHandler()))
}

// design and code by tsingson
