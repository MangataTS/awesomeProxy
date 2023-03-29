package main

import (
	"awesomeProxy/AsCache"
	"awesomeProxy/Core"
	"awesomeProxy/Log"
	"awesomeProxy/Reproxy"
	"awesomeProxy/Utils"
	"awesomeProxy/config"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func init() {
	//配置初始化
	config.CONFIG.Init()

	if config.CONFIG.ProxyMethod {
		config.CacheInit()
	} else {
		// 初始化并根证书
		err := Core.NewCertificate().Init()
		if err != nil {
			Log.Fatal("初始化根证书失败：" + err.Error())
		}
		//打开系统代理
		Host := "localhost" + config.CONFIG.CoProxy.Port
		Utils.SetWindowsProxy(Host)
	}

}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *AsCache.Group {
	return AsCache.NewGroup("scores", 2<<10, AsCache.GetterFunc(
		func(key string) ([]byte, error) {
			Log.Debug("[SlowDB] search key ", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func main() {
	//如果是进行反向代理代理
	if config.CONFIG.ProxyMethod {
		port := flag.String("port", config.CONFIG.ReProxy.Port, "listen port")
		flag.Parse()
		gee := createGroup()
		reverseUrl := fmt.Sprintf("http://%v:%d", config.Insts[0].Host, config.Insts[0].Port)
		remote, err := url.Parse(reverseUrl)
		if err != nil {
			panic(err)
		}
		Pproxy := Reproxy.GoReverseProxy(&Reproxy.RProxy{
			Remote: remote,
		})
		Pproxy.Set(config.Addrs...)
		gee.RegisterPeers(Pproxy)
		serveErr := http.ListenAndServe(":"+*port, Pproxy)
		if serveErr != nil {
			panic(serveErr)
		}
	} else {
		// 正常的正向网关代理
		port := flag.String("port", config.CONFIG.CoProxy.Port, "listen port")
		nagle := flag.Bool("nagle", config.CONFIG.CoProxy.Nagle, "connect remote use nagle algorithm")
		TcpProxy := flag.String("proxy", "", "prxoy remote host")
		to := flag.String("to", "", "tcp remote host")
		flag.Parse()
		if *port == "0" {
			Log.Fatal("port required")
			return
		}
		// 启动服务
		s := Core.NewProxyServer(*port, *nagle, *TcpProxy, *to)
		// 注册http客户端请求事件函数
		s.OnHttpRequestEvent = func(body []byte, request *http.Request, resolve Core.ResolveHttpRequest) {
			Log.Info("=========================HttpRequestEvent: =============================")
			Log.Info("Host : ", request.Host)
			Log.Info("Method : ", request.Method)
			Log.Info("URL : ", request.URL)
			Log.Info("Proto : ", request.Proto)
			Log.Info("Form : ", request.Form)
			Utils.BlacklistFilter(request)

			mimeType := request.Header.Get("Content-Type")
			if strings.Contains(mimeType, "json") {
				Log.Info("HttpRequestEvent：" + string(body))
			}
			// 可以在这里做数据修改
			resolve(body, request)
		}
		// 注册http服务器响应事件函数
		s.OnHttpResponseEvent = func(body []byte, response *http.Response, resolve Core.ResolveHttpResponse) {
			mimeType := response.Header.Get("Content-Type")
			if strings.Contains(mimeType, "json") {
				Log.Info("HttpResponseEvent：" + string(body))
			}
			// 可以在这里做数据修改
			resolve(body, response)
		}
		//s.OnHttpResponseEvent = func(response *http.Response) {
		//	contentType := response.Header.Get("Content-Type")
		//	var reader io.Reader
		//	if strings.Contains(contentType, "json") {
		//		reader = bufio.NewReader(response.Body)
		//		if header := response.Header.Get("Content-Encoding"); header == "gzip" {
		//			reader, _ = gzip.NewReader(response.Body)
		//		}
		//		body, _ := io.ReadAll(reader)
		//		Log.Info("HttpResponseEvent len:", len(string(body)))
		//		//Log.Log.Println("HttpResponseEvent：" + string(body))
		//	}
		//}
		// 注册socket5服务器推送消息事件函数
		s.OnSocket5ResponseEvent = func(message []byte, resolve Core.ResolveSocks5) (int, error) {
			Log.Info("Socket5ResponseEvent：" + string(message))
			// 可以在这里做数据修改
			return resolve(message)
		}
		// 注册socket5客户端推送消息事件函数
		s.OnSocket5RequestEvent = func(message []byte, resolve Core.ResolveSocks5) (int, error) {
			Log.Info("Socket5RequestEvent：" + string(message))
			// 可以在这里做数据修改
			return resolve(message)
		}
		// 注册ws客户端推送消息事件函数
		s.OnWsRequestEvent = func(msgType int, message []byte, resolve Core.ResolveWs) error {
			Log.Info("WsRequestEvent：" + string(message))
			// 可以在这里做数据修改
			return resolve(msgType, message)
		}
		// 注册ws服务器推送消息事件函数
		s.OnWsResponseEvent = func(msgType int, message []byte, resolve Core.ResolveWs) error {
			Log.Info("WsResponseEvent：" + string(message))
			// 可以在这里做数据修改
			return resolve(msgType, message)
		}

		// 注册w服务器推送消息事件函数
		s.OnTcpClientStreamEvent = func(message []byte, resolve Core.ResolveTcp) (int, error) {
			Log.Info("TcpClientStreamEvent：" + string(message))
			// 可以在这里做数据修改
			return resolve(message)
		}

		// 注册w服务器推送消息事件函数
		s.OnTcpServerStreamEvent = func(message []byte, resolve Core.ResolveTcp) (int, error) {
			Log.Info("TcpServerStreamEvent：" + string(message))
			// 可以在这里做数据修改
			return resolve(message)
		}
		_ = s.Start()
	}

}

/*
--port:代理服务监听的端口,默认为9090

--to:代理tcp服务时,目的服务器的ip和端口,默认为0,仅tcp代理使用

--proxy:上级代理地址

--nagle:是否开启nagle数据合并算法,默认true

*/
