package Reproxy

import (
	"awesomeProxy/Log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// CrawlerUserAgent 表示的是UA的爬虫正则匹配
var CrawlerUserAgent = []string{".*request.*", ".*python.*"}

// IpMap 用来存储每个ip的请求数和时间
var IpMap = SafeMap{m: make(map[string]int64)}

type SafeMap struct {
	m   map[string]int64
	mux sync.Mutex
}

func (sm *SafeMap) Inc(key string) {
	sm.mux.Lock()
	defer sm.mux.Unlock()
	sm.m[key]++
}

func (sm *SafeMap) Equal(key string, value int64) {
	sm.mux.Lock()
	defer sm.mux.Unlock()
	sm.m[key] = value
}

func (sm *SafeMap) Get(key string) (int64, bool) {
	sm.mux.Lock()
	defer sm.mux.Unlock()
	value, ok := sm.m[key]
	return value, ok
}

func (sm *SafeMap) DeleteIp(ip string) {
	if _, ok := IpMap.Get(ip); !ok {
		return
	}
	IpMap.mux.Lock()
	defer IpMap.mux.Unlock()
	delete(IpMap.m, ip)
}

// MaxIpQps 每分钟最大请求次数
var MaxIpQps = 60

// FiltrationCrawler 爬虫过滤函数，如果检测到是爬虫，那么就返回true，否则返回false
func FiltrationCrawler(request http.Request) bool {
	ip := getIP(request)
	Log.Info("GetIp :", ip)
	if UaCheck(request) {
		return true
	}
	// 正常请求，更新ip请求频率，如果是爬虫冷却时间，那么返回true
	return updateIPMap(ip)
}

// UaCheck 检查请求中的请求头是否正常
func UaCheck(request http.Request) bool {
	ua := request.UserAgent()
	ref := request.Referer()
	// 检测请求头长度 和refer地址
	if len(ua) < 80 && ref == "" {
		return true
	}
	for _, s := range CrawlerUserAgent {
		re, err := regexp.Compile(s)
		if err != nil {
			Log.Fatal("UaCheck 正则匹配出问题~")
		}
		if re.MatchString(ua) {
			return true
		}
	}
	return false
}

// 获取请求的ip地址，如果有多个，取第一个
func getIP(r http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	ips := strings.Split(ip, ", ")
	ipt := strings.Split(ips[0], ":")
	return ipt[0]
}

// 更新ipMap中对应的值，如果不存在则创建，如果存在则增加1，并记录当前时间戳（毫秒），并且检测爬虫上次爬取的时间
func updateIPMap(ip string) bool {
	now := time.Now().UnixNano() / 1e6 // 毫秒时间戳
	if _, ok := IpMap.Get(ip); !ok {
		// 新建一个键值对，初始请求数为1，初始时间为当前时间戳
		IpMap.Equal(ip, now*1000+1)
	} else {
		// 取出原来的键值对，分离出请求数和时间戳
		val, _ := IpMap.Get(ip)
		count := val % 1000 // 请求数在低三位
		ts := val / 1000    // 时间戳在高位
		// 如果当前时间戳和上次时间戳相差超过一分钟，则重置请求数为1，否则增加1
		if now-ts > 60*1000 {
			count = 1
		} else {
			// 发现ip的请求超过阈值
			if GetIpQps(ip) > int64(MaxIpQps) {
				return true
			}
			count++
		}
		// 更新键值对，重新组合请求数和时间戳
		IpMap.Equal(ip, now*1000+count)
	}
	return false
}

func GetIpQps(ip string) int64 {
	key, ok := IpMap.Get(ip)
	if !ok {
		return 0
	}
	return key % 1000
}
