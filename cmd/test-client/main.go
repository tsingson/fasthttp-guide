package main

import (
	"github.com/sanity-io/litter"

	"github.com/valyala/fasthttp"

	"github.com/tsingson/fasthttp-guide/webclient"

	"github.com/tsingson/fasthttp-guide/pkg/vtils"
)

func main() {
	w := webclient.Default()
	w.Debug = true

	w.Authentication = false
	w.JwtToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTY2NzIwMDAsInJvbGUiOiJ0ZXJtaW5hbF9hcGsiLCJzdGF0dXMiOiJhY3RpdmUiLCJ1c2VyX2lkIjoiNTBjNjg5MTAtNjEyYi00NjMzLTk2YjktNTA3NzhjNDViNTAwIn0.l1JHnOL85s3ajto0MKs-D6paW1YxpaMuxA0nzI0Xlfk"
	url := "http://localhost:3001/get"
	resp, err := w.FastGet(url)
	if err != nil {
	}
	if resp != nil {
		litter.Dump(vtils.B2S(resp.Body()))
	}
	// clean-up
	fasthttp.ReleaseResponse(resp)
	w.Authentication = false
	url = "http://localhost:3001/post"

	b := []byte(`{"actual_start_date":"2019-07-29","actual_end_date":"2019-07-29","plan_start_date":"2019-07-29","plan_end_date":"2019-02-12","title":"养殖计划00002","user_id":2098735545843717147}`)

	w.Accept = webclient.AcceptRest

	resp1, er1 := w.FastPostByte(url, b)

	if er1 != nil {
	}
	if resp1 != nil {
		litter.Dump(vtils.B2S(resp1.Body()))
	}
	// clenn-up
	fasthttp.ReleaseResponse(resp1)
}
