package Reproxy

import (
	"awesomeProxy/Log"
	"awesomeProxy/balance"
	"awesomeProxy/config"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type RProxy struct {
	Remote *url.URL
}

func GoReverseProxy(this *RProxy) *httputil.ReverseProxy {
	remote := this.Remote

	proxy := httputil.NewSingleHostReverseProxy(remote)

	proxy.Director = func(request *http.Request) {
		if FiltrationCrawler(*request) {
			Log.Warn("发现爬虫，现在正在对IP进行监管")
			request.URL.Scheme = "http"
			request.Host = "localhost"
			return
		}
		ins, err := balance.DoBalance(config.BalanceNames[config.CONFIG.ReProxy.BalanceMethod], config.Insts)
		if err != nil {
			Log.Fatal("report error" + err.Error())
		}
		Forw := fmt.Sprintf("http://%v:%d", ins.Host, ins.Port)
		remote, err = url.Parse(Forw)
		if err != nil {
			Log.Error("GoReverseProxy remote url" + err.Error())
		}
		targetQuery := remote.RawQuery
		request.URL.Scheme = remote.Scheme
		request.URL.Host = remote.Host
		request.Host = remote.Host // todo 这个是关键
		request.URL.Path, request.URL.RawPath = joinURLPath(remote, request.URL)

		if targetQuery == "" || request.URL.RawQuery == "" {
			request.URL.RawQuery = targetQuery + request.URL.RawQuery
		} else {
			request.URL.RawQuery = targetQuery + "&" + request.URL.RawQuery
		}
		//Log.Info("request.UserAgent() : ", request.UserAgent())
		Log.Info("request.URL.Path：", request.URL.Path, " request.URL.RawQuery：", request.URL.RawQuery)
	}

	// 修改响应头
	proxy.ModifyResponse = func(response *http.Response) error {
		response.Header.Add("Access-Control-Allow-Origin", "*")
		response.Header.Add("Reverse-Proxy-Server-PowerBy", "(kaptree)https://acm.mangata.ltd")
		return nil
	}

	return proxy
}

// go sdk 源码
func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

// go sdk 源码
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
