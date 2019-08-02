package webclient

import (
	"time"

	"github.com/savsgio/gotils"
	"github.com/tsingson/zaplogger"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

const (
	AcceptJson = "application/json"
	AcceptRest = "application/vnd.pgrst.object+json"
)

//
// authentication, authorization, and accounting

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

// FastPostByte  do  POST request via fasthttp
func (w *WebClient) FastPostByte(requestURI string, body []byte) (*fasthttp.Response, error) {
	var log = w.log.Named("FastPostByte")
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

	req.Header.SetMethod("POST")
	req.SetBody(body)
	// fmt.Println("---------- req --------------")
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
	var err = fasthttp.DoTimeout(req, resp, timeOut)
	if err != nil {
		log.Error("post request error", zap.Error(err))
		return nil, err
	}
	// list all response for debug
	if w.Debug {
		elapsed := time.Since(t1)
		log.Debug(w.TransactionID, zap.Duration("elapsed", elapsed))
		log.Debug(w.TransactionID, zap.Int("http status code", resp.StatusCode()))
		resp.Header.VisitAll(func(key, value []byte) {
			log.Debug(w.TransactionID, zap.String("key", gotils.B2S(key)), zap.String("value", gotils.B2S(value)))
		})
		log.Debug(w.TransactionID, zap.String("http payload", gotils.B2S(resp.Body())))
	}

	// just for demo
	var out = fasthttp.AcquireResponse()
	resp.CopyTo(out)

	return out, nil
}

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
