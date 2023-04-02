package Reproxy

import (
	"awesomeProxy/Log"
	"awesomeProxy/balance"
	"awesomeProxy/config"
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type RProxy struct {
	Remote *url.URL
}

func FiltIp(ip string) bool {
	file, err := os.Open("Forbid_IP.txt")
	if err != nil {
		Log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			Log.Error("文件关闭失败:", err)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	// 循环扫描每一行
	for scanner.Scan() {
		// 获取当前行的文本
		line := scanner.Text()
		// 打印当前行的文本
		if strings.Contains(ip, line) {
			Log.Warn("检测到黑名单IP: ", ip)
			return true
		}
	}
	// 检查扫描是否有错误
	if err := scanner.Err(); err != nil {
		Log.Fatal(err)
	}
	return false
}

func GoReverseProxy(this *RProxy) *ReverseProxy {
	remote := this.Remote

	proxy := NewSingleHostReverseProxy(remote)

	proxy.Director = func(request *http.Request) {

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
		request.Host = remote.Host
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
