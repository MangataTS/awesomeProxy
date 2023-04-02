package Core

import (
	"awesomeProxy/Contract"
	"awesomeProxy/Log"
	"awesomeProxy/Utils"
	"awesomeProxy/config"
	"bufio"
	"fmt"
	"github.com/viki-org/dnscache"
	"net"
	"net/http"
	"runtime"
	"time"
)

type HttpRequestEvent func(message []byte, request *http.Request, resolve ResolveHttpRequest)
type HttpResponseEvent func(message []byte, response *http.Response, resolve ResolveHttpResponse)

type Socket5ResponseEvent func(message []byte, resolve ResolveSocks5) (int, error)
type Socket5RequestEvent func(message []byte, resolve ResolveSocks5) (int, error)

type WsRequestEvent func(msgType int, message []byte, resolve ResolveWs) error
type WsResponseEvent func(msgType int, message []byte, resolve ResolveWs) error

type TcpServerStreamEvent func(message []byte, resolve ResolveTcp) (int, error)
type TcpClientStreamEvent func(message []byte, resolve ResolveTcp) (int, error)

const (
	MethodGet     = 0x47
	MethodConnect = 0x43
	MethodPost    = 0x50
	MethodPut     = 0x50
	MethodDelete  = 0x44
	MethodOptions = 0x4F
	MethodHead    = 0x48

	SocksFive = 0x5
)

type ProxyServer struct {
	nagle                  bool
	to                     string
	proxy                  string
	port                   string
	listener               *net.TCPListener
	dns                    *dnscache.Resolver
	OnHttpRequestEvent     HttpRequestEvent
	OnHttpResponseEvent    HttpResponseEvent
	OnWsRequestEvent       WsRequestEvent
	OnWsResponseEvent      WsResponseEvent
	OnSocket5ResponseEvent Socket5ResponseEvent
	OnSocket5RequestEvent  Socket5RequestEvent
	OnTcpServerStreamEvent TcpServerStreamEvent
	OnTcpClientStreamEvent TcpClientStreamEvent
}

func (i *ProxyServer) Install() {
	if runtime.GOOS == "windows" {
		err := Utils.InstallCert("cert.crt")
		if err != nil {
			Log.Fatal(err.Error())
		}
		Log.Info("已安装系统证书")
		err = Utils.SetWindowsProxy(fmt.Sprintf("localhost:%s", i.port))
		if err != nil {
			Log.Fatal(err.Error())
		}
		Log.Info("已设置系统代理")
	}
	Log.Info("非windows系统请手动安装证书并设置代理,可以在根目录或访问http://awesome-proxy.io/tls获取证书文件")
}

func NewProxyServer(port string, nagle bool, proxy string, to string) *ProxyServer {
	return &ProxyServer{
		port:  port,
		dns:   dnscache.New(time.Minute * 5),
		nagle: nagle,
		proxy: proxy,
	}
}

func (i *ProxyServer) Start() error {
	// 解析地址
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", i.port))
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	// 监听服务
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	i.listener = listener
	i.Install()
	i.MultiListen()
	select {}
}

// 这里协程嵌套处理了客户端请求处理，通过五个协程进行监听客户端操作，然后协程内部再通过协程handle处理conn链接

func (i *ProxyServer) MultiListen() {
	LoopNum := config.CONFIG.CoProxy.MultiListenNum
	for s := 0; s < LoopNum; s++ {
		go func() {
			for {
				conn, err := i.listener.Accept()
				if err != nil {
					if e, ok := err.(net.Error); ok && e.Timeout() {
						Log.Error("接受连接超时：" + err.Error())
						time.Sleep(time.Second / 20)
					} else {
						Log.Error("接受连接失败：" + err.Error())
					}
					continue
				}
				go i.handle(conn)
			}
		}()
	}
}

func (i *ProxyServer) handle(conn net.Conn) {
	var process Contract.IServerProcesser
	defer conn.Close()
	// 使用bufio读取,原conn的句柄数据被读完
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	// 预读取一段字节,https、ws、wss读取到的数据为：CONNECT xx.com:8080 HTTP/1.1
	peek, err := reader.Peek(3)
	if err != nil {
		return
	}
	peer := ConnPeer{server: i, conn: conn, writer: writer, reader: reader}
	switch peek[0] {
	case MethodGet, MethodPost, MethodDelete, MethodOptions, MethodHead, MethodConnect:
		process = &ProxyHttp{ConnPeer: peer}
		break
	case SocksFive:
		process = &ProxySocks5{ConnPeer: peer}
	default:
		process = &ProxyTcp{ConnPeer: peer}
	}
	process.Handle()
}
