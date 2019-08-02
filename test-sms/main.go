package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/tsingson/uuid"

	"github.com/sanity-io/litter"
	"github.com/savsgio/gotils"

	"github.com/tsingson/fasthttp-example/webclient"
)

func main() {

	var w = webclient.Default()
	w.Debug = true

	w.Authentication = false

	var resp, err = w.FastGet(buildURL())

	if err != nil {

	}
	if resp != nil {
		litter.Dump(gotils.B2S(resp.Body()))
	}
	// clean-up
	// fasthttp.ReleaseResponse(resp)

}

const sortQueryString_fmt string = "AccessKeyId=%s" +
	"&Action=SendSms" +
	"&Format=JSON" +
	"&OutId=123" +
	"&PhoneNumbers=%s" +
	"&RegionId=cn-hangzhou" +
	"&SignName=%s" +
	"&SignatureMethod=HMAC-SHA1" +
	"&SignatureNonce=%s" +
	"&SignatureVersion=1.0" +
	"&TemplateCode=%s" +
	"&TemplateParam=%s" +
	"&Timestamp=%s" +
	"&Version=2017-05-25"

func encode_local(encode_str string) string {
	urlencode := url.QueryEscape(encode_str)
	urlencode = strings.Replace(urlencode, "+", "%%20", -1)
	urlencode = strings.Replace(urlencode, "*", "%2A", -1)
	urlencode = strings.Replace(urlencode, "%%7E", "~", -1)
	urlencode = strings.Replace(urlencode, "/", "%%2F", -1)
	return urlencode
}

func buildURL() string {
	const token string = "GqgLZwEJY1Lu1H4RR2Q7ZNnEzXCabY&" // 阿里云 accessSecret 注意这个地方要添加一个 &

	AccessKeyId := "LTAIYYZkfMOWof4h" // 自己的阿里云 accessKeyID
	PhoneNumbers := "13602658546"     // 发送目标的手机号
	SignName := url.QueryEscape("阿里云短信测试专用")
	SignatureNonce, _ := uuid.NewV4()
	TemplateCode := "SMS_2945522"
	TemplateParam := url.QueryEscape(`{"code":"123456"","product":"test"}`)
	Timestamp := url.QueryEscape(time.Now().UTC().Format("2006-01-02T15:04:05Z"))

	sortQueryString := fmt.Sprintf(sortQueryString_fmt,
		AccessKeyId,
		PhoneNumbers,
		SignName,
		SignatureNonce,
		TemplateCode,
		TemplateParam,
		Timestamp,
	)

	urlencode := encode_local(sortQueryString)
	sign_str := fmt.Sprintf("GET&%%2F&%s", urlencode)

	key := []byte(token)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(sign_str))
	signture := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	signture = encode_local(signture)
	fmt.Println(signture)
	return fmt.Sprintf("http://dysmsapi.aliyuncs.com/?Signature=%s&%s\n", signture, sortQueryString)
}
