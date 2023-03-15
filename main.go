package main

import (
	"awesomeProxy/Core"
	"awesomeProxy/Core/Websocket"
	"awesomeProxy/Log"
	"awesomeProxy/Reproxy"
	"awesomeProxy/config"
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	// 初始化日志
	Log.NewLogger().Init()
	// 初始化根证书
	err := Core.NewCertificate().Init()
	if err != nil {
		Log.Log.Println("初始化根证书失败：" + err.Error())
		return
	}
	config.CONFIG.Init()
}

func main() {

	//如果是进行反向代理代理
	if config.CONFIG.ProxyMethod {
		port := flag.String("port", config.CONFIG.ReProxy.Port, "listen port")
		flag.Parse()
		reverseUrl := fmt.Sprintf("http://%v:%d", config.Insts[0].Host, config.Insts[0].Port)
		remote, err := url.Parse(reverseUrl)
		if err != nil {
			panic(err)
		}
		Pproxy := Reproxy.GoReverseProxy(&Reproxy.RProxy{
			Remote: remote,
		})
		log.Println("当前代理地址： " + reverseUrl + " 本地监听： http://127.0.0.1:" + *port)

		serveErr := http.ListenAndServe(":"+*port, Pproxy)

		if serveErr != nil {
			panic(serveErr)
		}
	} else {
		// 正常的正向网关代理
		port := flag.String("port", config.CONFIG.CoProxy.Port, "listen port")
		nagle := flag.Bool("nagle", config.CONFIG.CoProxy.Nagle, "connect remote use nagle algorithm")
		proxy := flag.String("proxy", "0", "tcp prxoy remote host")
		flag.Parse()
		if *port == "0" {
			Log.Log.Fatal("port required")
			return
		}
		// 启动服务
		s := Core.NewProxyServer(*port, *nagle, *proxy)
		// 注册http客户端请求事件函数
		s.OnHttpRequestEvent = func(request *http.Request) {
			Log.Log.Println("=========================HttpRequestEvent: =============================")
			Log.Log.Println("Host : ", request.Host)
			Log.Log.Println("Method : ", request.Method)
			Log.Log.Println("URL : ", request.URL)
			Log.Log.Println("Proto : ", request.Proto)
			Log.Log.Println("Form : ", request.Form)
			if request.Host == "127.0.0.1:10001" {
				return
			}
		}
		// 注册http服务器响应事件函数
		s.OnHttpResponseEvent = func(response *http.Response) {
			contentType := response.Header.Get("Content-Type")
			var reader io.Reader
			if strings.Contains(contentType, "json") {
				reader = bufio.NewReader(response.Body)
				if header := response.Header.Get("Content-Encoding"); header == "gzip" {
					reader, _ = gzip.NewReader(response.Body)
				}
				body, _ := io.ReadAll(reader)
				Log.Log.Println("HttpResponseEvent len:", len(string(body)))
				//Log.Log.Println("HttpResponseEvent：" + string(body))
			}
		}
		// 注册socket5服务器推送消息事件函数
		s.OnSocket5ResponseEvent = func(message []byte) {
			Log.Log.Println("Socket5ResponseEvent：" + string(message))
		}
		// 注册socket5客户端推送消息事件函数
		s.OnSocket5RequestEvent = func(message []byte) {
			Log.Log.Println("Socket5RequestEvent：" + string(message))
		}
		// 注册ws客户端推送消息事件函数
		s.OnWsRequestEvent = func(msgType int, message []byte, target *Websocket.Conn, resolve Core.ResolveWs) error {
			Log.Log.Println("WsRequestEvent：" + string(message))
			return target.WriteMessage(msgType, message)
		}
		// 注册ws服务器推送消息事件函数
		s.OnWsResponseEvent = func(msgType int, message []byte, client *Websocket.Conn, resolve Core.ResolveWs) error {
			Log.Log.Println("WsResponseEvent：" + string(message))
			return resolve(msgType, message, client)
		}
		_ = s.Start()
	}

}
