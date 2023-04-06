package main

import (
	"awesomeProxy/AsCache"
	"awesomeProxy/Core"
	"awesomeProxy/Log"
	"awesomeProxy/Report"
	"awesomeProxy/Reproxy"
	"awesomeProxy/Utils"
	"awesomeProxy/ac_automaton"
	"awesomeProxy/config"
	"awesomeProxy/global"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
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
		IsReport := flag.Bool("IsReport", false, "Reserve Server proxy Report")
		flag.Parse()
		if *IsReport {
			Report.GetReReport()
		}
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
		IsReport := flag.Bool("IsReport", false, "Reserve Server proxy Report")
		flag.Parse()
		if *IsReport {
			Report.GetCoReport()
		}
		if *port == "0" {
			Log.Fatal("port required")
			return
		}
		// 启动服务
		s := Core.NewProxyServer(*port, *nagle, *TcpProxy, *to)
		// 注册tcp连接事件
		s.OnTcpConnectEvent = func(conn net.Conn) {

		}
		// 注册tcp关闭事件
		s.OnTcpCloseEvent = func(conn net.Conn) {

		}

		// 注册http客户端请求事件函数
		s.OnHttpRequestEvent = func(body []byte, request *http.Request, resolve Core.ResolveHttpRequest, conn net.Conn) bool {
			Log.Debug("=========================HttpRequestEvent: =============================")
			Log.Debug("Host : ", request.Host)
			Log.Debug("Method : ", request.Method)
			Log.Debug("URL : ", request.URL)
			Log.Debug("Proto : ", request.Proto)
			Log.Debug("Form : ", request.Form)
			if Utils.BlacklistFilter(request) {
				// CalCoBlackHostData 数据统计 ok 上锁
				global.CalCoBlaLock.Lock()
				Uvalue, ok := global.CalCoBlackHostData[request.Host]
				if !ok {
					Uvalue = 0
				}
				Uvalue++
				global.CalCoBlackHostData[request.Host] = Uvalue
				global.CalCoBlaLock.Unlock()
				global.SaveCoConfig()
				return true
			}

			// CalCoRequestData 数据统计 上锁
			global.CalCoReqLock.Lock()
			value, ok := global.CalCoRequestData[request.Host]
			if !ok {
				value = 0
			}
			global.CalCoRequestData[request.Host] = value + 1
			global.CalCoReqLock.Unlock()
			global.SaveCoConfig()

			mimeType := request.Header.Get("Content-Type")
			if strings.Contains(mimeType, "json") {
				//Log.Info("HttpRequestEvent：" + string(body))
			}
			// 可以在这里做数据修改
			resolve(body, request)
			// 如果正常处理必须返回true，如果不需要发送请求，返回false，一般在自己操作conn的时候才会用到
			return true
		}

		// 注册http服务器响应事件函数
		s.OnHttpResponseEvent = func(body []byte, response *http.Response, resolve Core.ResolveHttpResponse, conn net.Conn) bool {
			mimeType := response.Header.Get("Content-Type")
			if strings.Contains(mimeType, "json") {
				Log.Info("HttpResponseEvent：" + string(body))
			}
			// todo 对其他协议进行数据统计
			if response.TLS != nil {
				// HTTPS 数据统计 上锁
				global.CalCoProLock.Lock()
				tvalue, ok := global.CalCoProtocolData["HTTPS"]
				if !ok {
					tvalue.Name = "HTTPS"
					global.Glock.Lock()
					tvalue.ReqTimes = global.CoReportConfig.CoProtocolData[1].ReqTimes
					tvalue.ReqDataSize = global.CoReportConfig.CoProtocolData[1].ReqDataSize
					global.Glock.Unlock()
				}
				tvalue.ReqTimes++
				tvalue.ReqDataSize += len(body)
				global.CalCoProtocolData["HTTPS"] = tvalue
				global.CalCoProLock.Unlock()
				global.SaveCoConfig()
			} else {
				// HTTP 数据统计 上锁
				global.CalCoProLock.Lock()
				tvalue, ok := global.CalCoProtocolData["HTTP"]
				if !ok {
					tvalue.Name = "HTTP"
					global.Glock.Lock()
					tvalue.ReqTimes = global.CoReportConfig.CoProtocolData[0].ReqTimes
					tvalue.ReqDataSize = global.CoReportConfig.CoProtocolData[0].ReqDataSize
					global.Glock.Unlock()
				}
				tvalue.ReqTimes++
				tvalue.ReqDataSize += len(body)
				global.CalCoProtocolData["HTTP"] = tvalue
				global.CalCoProLock.Unlock()
				global.SaveCoConfig()
			}
			sbody := string(body)
			res := ac_automaton.Acauto.FindMatches(sbody)
			mingcnt := 0
			tiggercnt := 0
			if len(res) != 0 {
				for name, cnt := range res {
					//Log.Debug("敏感词: ", name, " 触发次数: ", cnt, "拦截")
					mingcnt++
					tiggercnt += cnt
					sbody = strings.ReplaceAll(sbody, name, "*")
				}
				body = []byte(sbody) //[]byte("触发敏感词")
				Log.Info("敏感词数量：", mingcnt, " 触发总次数： ", tiggercnt)

				// 更新敏感网站 ok

				global.Glock.Lock()
				global.CoReportConfig.CoSensitiveData.TriggerNum += mingcnt
				global.CoReportConfig.CoSensitiveData.Interceptions += tiggercnt
				global.Glock.Unlock()
				global.CalCoSenLock.Lock()
				global.CalCoSensitiveDataUrl[response.Request.Host] = true
				global.CalCoSenLock.Unlock()
				global.SaveCoConfig()
			}
			// 可以在这里做数据修改
			resolve(body, response)
			// 如果正常处理必须返回true，如果不需要发送请求，返回false，一般在自己操作conn的时候才会用到
			return true
		}
		// 注册socket5服务器推送消息事件函数
		s.OnSocks5ResponseEvent = func(message []byte, resolve Core.ResolveSocks5, conn net.Conn) (int, error) {
			Log.Info("Socks5ResponseEvent：" + string(message))
			// SOCKS5 数据统计 上锁
			global.CalCoProLock.Lock()
			tvalue, ok := global.CalCoProtocolData["SOCKS5"]
			if !ok {
				tvalue.Name = "SOCKS5"
				global.Glock.Lock()
				tvalue.ReqTimes = global.CoReportConfig.CoProtocolData[2].ReqTimes
				tvalue.ReqDataSize = global.CoReportConfig.CoProtocolData[2].ReqDataSize
				global.Glock.Unlock()
			}
			tvalue.ReqTimes++
			tvalue.ReqDataSize += len(message)
			global.CalCoProtocolData["SOCKS5"] = tvalue
			global.CalCoProLock.Unlock()
			global.SaveCoConfig()

			// 可以在这里做数据修改
			return resolve(message)
		}
		// 注册socket5客户端推送消息事件函数
		s.OnSocks5RequestEvent = func(message []byte, resolve Core.ResolveSocks5, conn net.Conn) (int, error) {
			Log.Info("Socks5RequestEvent：" + string(message))
			// 可以在这里做数据修改
			return resolve(message)
		}
		// 注册ws客户端推送消息事件函数
		s.OnWsRequestEvent = func(msgType int, message []byte, resolve Core.ResolveWs, conn net.Conn) error {
			Log.Info("WsRequestEvent：" + string(message))
			// 可以在这里做数据修改
			return resolve(msgType, message)
		}
		// 注册ws服务器推送消息事件函数
		s.OnWsResponseEvent = func(msgType int, message []byte, resolve Core.ResolveWs, conn net.Conn) error {
			Log.Info("WsResponseEvent：" + string(message))
			if conn.(*tls.Conn) != nil {
				// WebSocket TLS 数据统计 上锁
				global.CalCoProLock.Lock()
				tvalue, ok := global.CalCoProtocolData["WS"]
				if !ok {
					tvalue.Name = "WSS"
					global.Glock.Lock()
					tvalue.ReqTimes = global.CoReportConfig.CoProtocolData[4].ReqTimes
					tvalue.ReqDataSize = global.CoReportConfig.CoProtocolData[4].ReqDataSize
					global.Glock.Unlock()
				}
				tvalue.ReqTimes++
				tvalue.ReqDataSize += len(message)
				global.CalCoProtocolData["WSS"] = tvalue
				global.CalCoProLock.Unlock()
				global.SaveCoConfig()
			} else {
				// WebSocket 数据统计 上锁
				global.CalCoProLock.Lock()
				tvalue, ok := global.CalCoProtocolData["WS"]
				if !ok {
					tvalue.Name = "WS"
					global.Glock.Lock()
					tvalue.ReqTimes = global.CoReportConfig.CoProtocolData[3].ReqTimes
					tvalue.ReqDataSize = global.CoReportConfig.CoProtocolData[3].ReqDataSize
					global.Glock.Unlock()
				}
				tvalue.ReqTimes++
				tvalue.ReqDataSize += len(message)
				global.CalCoProtocolData["WS"] = tvalue
				global.CalCoProLock.Unlock()
				global.SaveCoConfig()
			}
			// 可以在这里做数据修改
			return resolve(msgType, message)
		}

		// 注册tcp服务器推送消息事件函数
		s.OnTcpClientStreamEvent = func(message []byte, resolve Core.ResolveTcp, conn net.Conn) (int, error) {
			Log.Info("TcpClientStreamEvent：" + string(message))
			// 可以在这里做数据修改
			return resolve(message)
		}

		// 注册tcp服务器推送消息事件函数
		s.OnTcpServerStreamEvent = func(message []byte, resolve Core.ResolveTcp, conn net.Conn) (int, error) {
			Log.Info("TcpServerStreamEvent：" + string(message))
			//  TCP 数据统计 上锁
			global.CalCoProLock.Lock()
			tvalue, ok := global.CalCoProtocolData["TCP"]
			if !ok {
				tvalue.Name = "TCP"
				global.Glock.Lock()
				tvalue.ReqTimes = global.CoReportConfig.CoProtocolData[5].ReqTimes
				tvalue.ReqDataSize = global.CoReportConfig.CoProtocolData[5].ReqDataSize
				global.Glock.Unlock()
			}
			tvalue.ReqTimes++
			tvalue.ReqDataSize += len(message)
			global.CalCoProtocolData["TCP"] = tvalue
			global.CalCoProLock.Unlock()
			global.SaveCoConfig()

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
