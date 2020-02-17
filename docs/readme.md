# FastHTTP 开发指南 ( Fasthttp Programming Guide )

## 1. 动机 Motivation

近3年来, fasthttp 被我用在几个重大项目(对我而言, 项目有多重大, 与收钱的多少成正比) 中,  FastHTTP 在生产环境中应用积累了一些经验, 写了一些小文章介绍 FastHTTP 的开发.

考虑到 FastHTTP 的开发思路与 Golang 标准 net/http 库有很大不同,  索性写成一本小册子, 汇集众多 FastHTTP 的优秀文章, 介绍 fasthttp 的实现思路, 个人实际使用与经验得失.

 

In the past 3 years, I have used fasthttp in several major projects (for me, how important the project is, which is proportional to the amount of money collected), FastHTTP has accumulated some experience in the production environment, and wrote some small articles Introduce the development of FastHTTP.

Considering that the development idea of FastHTTP is very different from the Golang standard net / http library, I simply wrote a booklet that brings together many excellent articles of FastHTTP to introduce the implementation ideas of fasthttp, personal actual use and experience gains and losses.

this booklet is inspired by [https://github.com/astaxie/build-web-application-with-golang](https://github.com/astaxie/build-web-application-with-golang)

## 2. 关于 fasthttp 的优点介绍

关于 FastHTTP 
[简述]  [github.com/valyala/fasthttp](https://github.com/valyala/fasthttp) 是 golang 中一个标志性的高性能 HTTP库, 主要用于 webserver 开发, 以及 web client / proxy 等. fasthttp 的高性能开发思路, 启发了很多开发者.



fasthttp 自己的介绍如下:

> Fast HTTP package for Go. Tuned for high performance. Zero memory allocations in hot paths. Up to 10x faster than net/http
>
> Fast HTTP implementation for Go.
>
>
> Currently fasthttp is successfully used by [VertaMedia](https://vertamedia.com/)
> in a production serving up to 200K rps from more than 1.5M concurrent keep-alive
> connections per physical server.


事实上, 这有点小夸张, 但在一定场景下经过优化部署, 确是有很高的性能. 


以下文字来自 [傅小黑](https://my.oschina.net/fuxiaohei) 原创文章: [Go 开发 HTTP 的另一个选择 fasthttp](https://my.oschina.net/fuxiaohei/blog/753977) 写于2016/09/30 :

> fasthttp 是 Go 的一款不同于标准库 net/http 的 HTTP 实现。fasthttp 的性能可以达到标准库的 10 倍，说明他魔性的实现方式。主要的点在于四个方面：
>
> * net/http 的实现是一个连接新建一个 goroutine；fasthttp 是利用一个 worker 复用 goroutine，减轻 runtime 调度 goroutine 的压力
> * net/http 解析的请求数据很多放在 map[string]string(http.Header) 或 map[string][]string(http.Request.Form)，有不必要的 []byte 到 string 的转换，是可以规避的
> * net/http 解析 HTTP 请求每次生成新的 *http.Request 和 http.ResponseWriter; fasthttp 解析 HTTP 数据到 *fasthttp.RequestCtx，然后使用 sync.Pool 复用结构实例，减少对象的数量
> * fasthttp 会延迟解析 HTTP 请求中的数据，尤其是 Body 部分。这样节省了很多不直接操作 Body 的情况的消耗
>
> 但是因为 fasthttp 的实现与标准库差距较大，所以 API 的设计完全不同。使用时既需要理解 HTTP 的处理过程，又需要注意和标准库的差别。

这段文字非常精练的总结了 fasthttp 的特点, 我摘录了这部分放在这里, 感谢 [傅小黑](https://my.oschina.net/fuxiaohei)  --- 另外, [傅小黑](https://my.oschina.net/fuxiaohei) 的技术文章非常棒, 欢迎大家去围观他....



关于 FastHTTP 的性能比较, 参考如下:

* [谁是最快的 web 框架]( https://colobu.com/2016/04/06/the-fastest-golang-web-framework/)  2016 [colobu](http://weibo.com/colobu/) [website](https://colobu.com/about/)



## 3. 关于本指南

本指南以 web 应用( web application ) 开发为主线, 化繁为易,分层渐进, 简洁直接,注重实效.   step by step 分享 golang 在 web 应用开发及周边应用开发中的基础实践

本指南每一个章节, 均有示例代码示范, 大部分文字均以代码 + 文本 + 手绘配图为主要形式, 力求清晰准确说明概念/思路/代码实现 及相关技巧, 同时附上相关参考索引

本指南分为以下几大部分:
1. go 开发环境, 开发工具配置及简要的开发/编译/运行/调试基础
2. HTTP 1.1 基础, 从基本 TCP client/server 实例到实现一个简单的 HTTP client / server 
3. FastHTTP 基本应用实例
4. 

## 4. 关于我


网名 tsingson (三明智), 现居中国深圳.南山

原 ustarcom IPTV/OTT 事业部播控产品线技术架构湿/解决方案工程湿角色(8年), 自由职业者,

喜欢音乐(口琴,是第三/四/五届广东国际口琴嘉年华的主策划人之一), 摄影与越野, 

喜欢 golang 语言 (商用项目中主要用 postgres + golang )  

[https://github.com/tsingson](https://github.com/tsingson)
[Email: tsingson@me.com](mailto:tsingson@me.com)

2020/02/10 中国深圳.小罗号口琴音乐中心