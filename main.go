package main

import (
	"awesomeProxy/AsCache"
	"awesomeProxy/Core"
	"awesomeProxy/Core/Websocket"
	"awesomeProxy/Log"
	"awesomeProxy/Reproxy"
	"awesomeProxy/Utils"
	"awesomeProxy/config"
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	//配置初始化
	config.CONFIG.Init()
	// 初始化并根证书
	err := Core.NewCertificate().Init()
	if err != nil {
		Log.Fatal("初始化根证书失败：" + err.Error())
	}
	//打开系统代理
	if !config.CONFIG.ProxyMethod {
		Host := "127.0.0.1:" + config.CONFIG.CoProxy.Port
		Utils.SetWindowsProxy(Host)
	}

}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	//如果是进行反向代理代理
	if config.CONFIG.ProxyMethod {
		// 分布式缓存组，下面是测试代码
		// curl http://localhost:9090/ascache/scores/Tom
		// curl http://localhost:9090/ascache/scores/kkk
		// 编写回调Getter函数，当缓存不存在就在这里获取，后续启用可将其写为SQL语句或者其他方法
		AsCache.NewGroup("scores", 2<<10, AsCache.GetterFunc(
			func(key string) ([]byte, error) {
				Log.Debug("[SlowDB] search key", key)
				if v, ok := db[key]; ok {
					return []byte(v), nil
				}
				return nil, fmt.Errorf("%s not exist", key)
			}))
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
		serveErr := http.ListenAndServe(":"+*port, Pproxy)
		if serveErr != nil {
			panic(serveErr)
		}
	} else {
		// 正常的正向网关代理
		port := flag.String("port", config.CONFIG.CoProxy.Port, "listen port")
		nagle := flag.Bool("nagle", config.CONFIG.CoProxy.Nagle, "connect remote use nagle algorithm")
		TcpProxy := flag.String("proxy", "0", "tcp prxoy remote host")
		flag.Parse()
		if *port == "0" {
			Log.Fatal("port required")
			return
		}
		// 启动服务
		s := Core.NewProxyServer(*port, *nagle, *TcpProxy)
		// 注册http客户端请求事件函数
		s.OnHttpRequestEvent = func(request *http.Request) {
			Log.Info("=========================HttpRequestEvent: =============================")
			Log.Info("Host : ", request.Host)
			Log.Info("Method : ", request.Method)
			Log.Info("URL : ", request.URL)
			Log.Info("Proto : ", request.Proto)
			Log.Info("Form : ", request.Form)
			Utils.BlacklistFilter(request)
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
				Log.Info("HttpResponseEvent len:", len(string(body)))
				//Log.Log.Println("HttpResponseEvent：" + string(body))
			}
		}
		// 注册socket5服务器推送消息事件函数
		s.OnSocket5ResponseEvent = func(message []byte) {
			Log.Info("Socket5ResponseEvent：" + string(message))
		}
		// 注册socket5客户端推送消息事件函数
		s.OnSocket5RequestEvent = func(message []byte) {
			Log.Info("Socket5RequestEvent：" + string(message))
		}
		// 注册ws客户端推送消息事件函数
		s.OnWsRequestEvent = func(msgType int, message []byte, target *Websocket.Conn, resolve Core.ResolveWs) error {
			Log.Info("WsRequestEvent：" + string(message))
			return target.WriteMessage(msgType, message)
		}
		// 注册ws服务器推送消息事件函数
		s.OnWsResponseEvent = func(msgType int, message []byte, client *Websocket.Conn, resolve Core.ResolveWs) error {
			Log.Info("WsResponseEvent：" + string(message))
			return resolve(msgType, message, client)
		}
		_ = s.Start()
	}

}
