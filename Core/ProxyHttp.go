package Core

import (
	"awesomeProxy/Core/Websocket"
	"awesomeProxy/Log"
	"awesomeProxy/Utils"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const ConnectSuccess = "HTTP/1.1 200 Connection Established\r\n\r\n"
const ConnectFailed = "HTTP/1.1 502 Bad Gateway\r\n\r\n"
const SslFileHost = "awesome-proxy.io"

type ProxyHttp struct {
	ConnPeer
	request  *http.Request
	response *http.Response
	upgrade  *Websocket.Upgrader
	target   net.Conn
	tls      bool
	port     string
}

type ResolveHttpRequest func(message []byte, request *http.Request)
type ResolveHttpResponse func(message []byte, response *http.Response)

type ResolveWs func(msgType int, message []byte) error

// Handle tcp连接处理入口
func (i *ProxyHttp) Handle() {
	request, err := http.ReadRequest(i.reader)
	if err != nil {
		return
	}
	i.port = "-1"
	if hostname := strings.Split(request.Host, ":"); len(hostname) > 1 {
		i.port = hostname[len(hostname)-1]
	}
	i.request = request
	//CONNECT方法用来建立到给定URI标识的服务器的隧道；它通过简单的TCP / IP隧道更改请求连接，通常实使用解码的HTTP代理来进行SSL编码的通信（HTTPS）。
	// 如果是connect方法则是https请求或者ws、wss请求
	if i.request.Method == http.MethodConnect {
		i.tls = true
		i.handleSslRequest()
		return
	}
	// 否则是普通请求
	i.tls = false
	i.handleRequest()
}

// handleRequest 处理Http请求入口
func (i *ProxyHttp) handleRequest() {
	var err error
	if i.request.URL == nil {
		Log.Error("请求地址为空")
		return
	}
	// 如果是下载证书,返回证书，给HTTPS提供加密操作的
	if i.request.URL.Path == "/tls" {
		response := http.Response{
			StatusCode: http.StatusOK,
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header: http.Header{
				"Content-Type":              []string{"application/x-x509-ca-cert"},
				"Content-Disposition":       []string{"attachment;filename=cert.crt"},
				"Content-Transfer-Encoding": []string{"binary"},
			},
			Body: io.NopCloser(bytes.NewReader(Cert.RootCaStr)),
		}
		_ = response.Write(i.conn)
		return
	}
	resolveRequest := ResolveHttpRequest(func(message []byte, request *http.Request) {
		request.Body = io.NopCloser(bytes.NewReader(message))
		request.Header.Set("Content-Length", strconv.Itoa(len(message)))

	})
	body, _ := i.ReadRequestBody(i.request.Body)
	resolveResult := i.server.OnHttpRequestEvent(body, i.request, resolveRequest, i.conn)
	if !resolveResult {
		return
	}
	// 处理正常请求,获取响应，将客户端数据转发给请求的服务器
	i.response, err = i.Transport(i.request)
	if i.response == nil {
		Log.Error("远程服务器无响应-1")
		return
	}
	if err != nil {
		Log.Error("获取远程服务器响应失败：" + err.Error())
		return
	}
	body, _ = i.ReadResponseBody(i.response)
	resolveResponse := ResolveHttpResponse(func(message []byte, response *http.Response) {
		response.Body = io.NopCloser(bytes.NewReader(message))
		// 手动计算长度
		response.Header.Set("Content-Length", strconv.Itoa(len(message)))
	})
	resolveResult = i.server.OnHttpResponseEvent(body, i.response, resolveResponse, i.conn)
	if !resolveResult {
		return
	}
	_ = i.response.Write(i.conn)
}

// ReadRequestBody 读取http请求体
func (i *ProxyHttp) ReadRequestBody(reader io.Reader) ([]byte, error) {
	if reader == nil {
		return []byte{}, nil
	}
	body, err := io.ReadAll(reader)
	return body, err
}

func (i *ProxyHttp) ReadResponseBody(response *http.Response) ([]byte, error) {
	var reader io.Reader
	var err error
	reader = bufio.NewReader(response.Body)
	if header := response.Header.Get("Content-Encoding"); header == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return []byte{}, nil
		}
	}
	return io.ReadAll(reader)
}

// RemoveHeader 移除请求头
func (i *ProxyHttp) RemoveHeader(header http.Header) {
	removeHeaders := []string{
		"Keep-Alive",
		"Transfer-Encoding",
		"TE",
		"Connection",
		"Trailer",
		"Upgrade",
		"Proxy-Authorization",
		"Proxy-Authenticate",
		"Connection",
		"Accept-Encoding",
	}
	for _, value := range removeHeaders {
		if v := header.Get(value); len(v) > 0 {
			if strings.EqualFold(value, "Connection") {
				for _, customerHeader := range strings.Split(value, ",") {
					header.Del(strings.Trim(customerHeader, " "))
				}
			}
			header.Del(value)
		}
	}
}

// Transport http请求转发
// Transport源码详解：https://www.dandelioncloud.cn/article/details/1495052443728838658
func (i *ProxyHttp) Transport(request *http.Request) (*http.Response, error) {
	// 去除一些头部
	i.RemoveHeader(request.Header)
	transport := &http.Transport{
		DisableKeepAlives:     true,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
		DialContext:           i.DialContext(),
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
	}
	if i.ConnPeer.server.proxy != "" {
		transport.Proxy = http.ProxyURL(&url.URL{Host: i.server.proxy})
	}
	response, err := transport.RoundTrip(request)
	if err != nil {
		Log.Error("转发http请求失败：" + err.Error())
		return nil, err
	}
	// 去除一些头部
	i.RemoveHeader(response.Header)
	return response, err
}

// 处理tls请求
func (i *ProxyHttp) handleSslRequest() {
	var err error
	// 如果使用了上级代理
	if i.ConnPeer.server.proxy != "" {
		i.target, err = net.Dial("tcp", i.server.proxy)
	} else {
		if i.port == "443" {
			i.target, err = tls.Dial("tcp", i.request.Host, &tls.Config{
				InsecureSkipVerify: true,
			})
		} else {
			i.target, err = net.Dial("tcp", i.request.Host)
		}
	}

	if err != nil {
		_, err = i.conn.Write([]byte(ConnectFailed))
		return
	}
	_ = i.target.Close()
	// 向源连接返回连接成功
	_, err = i.conn.Write([]byte(ConnectSuccess))
	if err != nil {
		Log.Error("返回连接状态失败：" + err.Error())
		return
	}
	// 建立TLS连接并返回给源
	i.SslReceiveSend()
}

// SetRequest 设置请求头
func (i *ProxyHttp) SetRequest(request *http.Request) *http.Request {
	request.Header.Set("Connection", "false")
	request.URL.Host = request.Host
	request.URL.Scheme = "https"
	return request
}

// SslReceiveSend tls数据接收发送
func (i *ProxyHttp) SslReceiveSend() {
	var err error
	certificate, err := Cache.GetCertificate(i.request.Host, i.port)
	if err != nil {
		Log.Error(i.request.Host + "：获取证书失败：" + err.Error())
		return
	}
	if _, ok := certificate.(tls.Certificate); !ok {
		return
	}
	cert := certificate.(tls.Certificate)
	sslConn := tls.Server(i.conn, &tls.Config{
		Certificates: []tls.Certificate{cert},
	})
	// ssl校验
	err = sslConn.Handshake()
	// 如果不是http的TLS请求,则说明是普通ws请求(ws请求会TLS校验报错),这里专门处理这种情况
	if err != nil {
		i.tls = false
		if err == io.EOF || strings.Index(err.Error(), "closed") != -1 {
			Log.Error("客户端TLS握手失败: " + err.Error())
			return
		}
		// 反射读取最后一帧原始数据
		i.handleWsShakehandErr(Utils.GetLastTimeFrame(sslConn, "rawInput"))
		return
	}
	_ = sslConn.SetDeadline(time.Now().Add(time.Second * 60))
	i.conn = sslConn
	i.tls = true
	i.reader = bufio.NewReader(i.conn)
	i.request, err = http.ReadRequest(i.reader)
	if err != nil {
		if err == io.EOF {
			Log.Error("浏览器TLS连接断开：" + err.Error())
			return
		}
		Log.Error("读取TLS连接请求数据失败：" + err.Error())
		return
	}
	// 如果包含upgrade同步说明是wss请求
	if i.request.Header.Get("Connection") == "Upgrade" {
		i.handleWssRequest()
		return
	}
	resolveRequest := ResolveHttpRequest(func(message []byte, request *http.Request) {
		request.Body = io.NopCloser(bytes.NewReader(message))
		request.Header.Set("Content-Length", strconv.Itoa(len(message)))
	})

	i.request = i.SetRequest(i.request)
	body, _ := i.ReadRequestBody(i.request.Body)
	resolveResult := i.server.OnHttpRequestEvent(body, i.request, resolveRequest, i.conn)
	if !resolveResult {
		return
	}
	i.response, err = i.Transport(i.request)
	if err != nil {
		Log.Error("远程服务器响应失败：" + err.Error())
		return
	}
	if i.response == nil {
		Log.Error("远程服务器无响应-2")
		return
	}

	body, _ = i.ReadResponseBody(i.response)
	resolveResponse := ResolveHttpResponse(func(message []byte, response *http.Response) {
		response.Body = io.NopCloser(bytes.NewReader(message))
		// 手动计算长度
		response.Header.Set("Content-Length", strconv.Itoa(len(message)))
	})
	resolveResult = i.server.OnHttpResponseEvent(body, i.response, resolveResponse, i.conn)
	if !resolveResult {
		return
	}
	err = i.response.Write(i.conn)
	if err != nil {
		if strings.Contains(err.Error(), "aborted") {
			Log.Error("代理返回响应数据失败：连接已关闭")
			return
		}
		Log.Error("代理返回响应数据失败：" + err.Error())
	}
}

// 加密wss请求
func (i *ProxyHttp) handleWssRequest() {
	i.handleWsRequest()
}

// 普通ws请求
func (i *ProxyHttp) handleWsShakehandErr(rawProtolInput []byte) {
	var err error
	// 获取浏览器发送给服务器的头部和数据,构建一个完整的请求对象
	rawInput := string(rawProtolInput)
	_ = i.conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	rawInputList := strings.Split(rawInput, "\r\n")
	// 读取请求方法
	wsMethodList := strings.Split(rawInputList[0], " ")
	// 构建请求
	wsRequest := &http.Request{
		Method: wsMethodList[0],
		Header: map[string][]string{},
	}
	for _, value := range rawInputList {
		// 填充header
		headerKeValList := strings.Split(value, ": ")
		if len(headerKeValList) <= 1 {
			continue
		}
		wsRequest.Header.Set(headerKeValList[0], headerKeValList[1])
		// 填充host
		if headerKeValList[0] == "Host" {
			wsRequest.Host = headerKeValList[1]
			wsRequest.RequestURI = fmt.Sprintf("http://%s", wsRequest.Host)
			wsRequest.URL, err = url.Parse(fmt.Sprintf("%s%s", wsRequest.RequestURI, wsMethodList[1]))
			if err != nil {
				Log.Error("解析ws请求地址错误：" + err.Error())
				return
			}
		}
		// 计算长度
		if headerKeValList[0] == "Content-Length" {
			rawLen := len(rawInput)
			bodyLen, _ := strconv.Atoi(headerKeValList[1])
			headerLen := rawLen - bodyLen - 4
			wsRequest.Body = io.NopCloser(bytes.NewBuffer([]byte(rawInput[headerLen:])))
			wsRequest.ContentLength = int64(bodyLen)
		}
	}
	i.request = wsRequest
	i.handleWsRequest()
}

// 处理ws请求
func (i *ProxyHttp) handleWsRequest() bool {
	if i.request.Header.Get("Upgrade") == "" {
		return false
	}
	if i.upgrade == nil {
		i.upgrade = &Websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
	}

	// 如果有携带自定义参数,添加上
	i.upgrade.Subprotocols = []string{i.request.Header.Get("Sec-WebSocket-Protocol")}
	recorder := httptest.NewRecorder()
	// 开始ws校验,向浏览器返回Sec-WebSocket-Accept头,ws握手完成
	clientWsConn, err := i.upgrade.Upgrade(recorder, i.request, nil, i.conn, bufio.NewReadWriter(i.reader, i.writer))
	if err != nil {
		Log.Error("升级ws协议失败：" + err.Error())
		return true
	}
	hostname := fmt.Sprintf("%s://%s%s", func() string {
		if i.tls {
			return "wss"
		}
		return "ws"
	}(), i.request.Host, i.request.URL.Path)
	if i.request.URL.RawQuery != "" {
		hostname += "?" + i.request.URL.RawQuery
	}
	// 去掉ws的头部,因为后续工具类会自己生成并附加到请求中
	i.RemoveWsHeader()
	var dialer Websocket.Dialer
	dialer = Websocket.Dialer{}
	// 如果是wss,客户端传输层忽略证书校验
	if i.tls {
		dialer = Websocket.Dialer{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			HandshakeTimeout: time.Second * 10,
		}
	}
	dialer.NetDialContext = i.DialContext()
	targetWsConn, response, err := dialer.Dial(hostname, i.request.Header)
	if err != nil {
		var header []byte
		if response != nil {
			header, _ = httputil.DumpResponse(response, false)
		}
		Log.Error("连接ws服务器失败：" + string(header) + err.Error())
		return true
	}
	defer func() {
		_ = targetWsConn.Close()
	}()
	stop := make(chan error, 2)
	// 读取浏览器数据(长连接)
	go func() {
		for {
			msgType, message, err := targetWsConn.ReadMessage()
			if err != nil {
				if Websocket.IsUnexpectedCloseError(err, Websocket.CloseGoingAway, Websocket.CloseAbnormalClosure) {
					stop <- fmt.Errorf("读取ws服务器数据失败-1：%w", err)
					break
				}
				stop <- fmt.Errorf("读取ws服务器数据失败-2：%w", err)
				break
			}
			err = i.server.OnWsResponseEvent(msgType, message, func(msgType int, message []byte) error {
				return clientWsConn.WriteMessage(msgType, message)
			}, i.conn)
			if err != nil {
				stop <- fmt.Errorf("发送ws浏览器数据失败-1：%w", err)
			}
		}
	}()
	go func() {
		for {
			msgType, message, err := clientWsConn.ReadMessage()
			if err != nil {
				if Websocket.IsUnexpectedCloseError(err, Websocket.CloseGoingAway, Websocket.CloseAbnormalClosure) {
					stop <- fmt.Errorf("读取ws浏览器数据失败-1：%w", err)
					break
				}
				stop <- fmt.Errorf("读取ws浏览器数据失败-2：%w", err)
				break
			}
			err = i.server.OnWsRequestEvent(msgType, message, func(msgType int, message []byte) error {
				return targetWsConn.WriteMessage(msgType, message)
			}, i.conn)
			if err != nil {
				stop <- fmt.Errorf("发送ws浏览器数据失败-1：%w", err)
			}
		}
	}()
	err = <-stop
	Log.Error(err.Error())
	return false
}

func (i *ProxyHttp) DialContext() func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
	return func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
		separator := strings.LastIndex(addr, ":")
		ipList, err := i.server.dns.Fetch(addr[:separator])
		var ip string
		for _, item := range ipList {
			ip = item.String()
			if !strings.Contains(ip, ":") {
				break
			}
		}
		tcpAddr, _ := net.ResolveTCPAddr("tcp", ip+addr[separator:])
		conn, err = net.DialTimeout("tcp", tcpAddr.String(), time.Duration(10)*time.Second)
		if err != nil {
			return conn, err
		}
		// 是否关闭nagle算法
		tcpConn, ok := conn.(*net.TCPConn)
		if ok {
			_ = tcpConn.SetNoDelay(i.server.nagle)
		}
		return tcpConn, err
	}
}

// 连接是否可用
func (i *ProxyHttp) WsIsConnected(conn *Websocket.Conn) bool {
	err := conn.WriteMessage(1, nil)
	return err == nil
}

// 移除ws请求头
func (i *ProxyHttp) RemoveWsHeader() {
	headers := []string{
		"Upgrade",
		"Connection",
		"Sec-Websocket-Key",
		"Sec-Websocket-Version",
		"Sec-Websocket-Extensions",
	}
	for _, value := range headers {
		if ok := i.request.Header.Get(value); ok != "" {
			i.request.Header.Del(value)
		}
	}
}
