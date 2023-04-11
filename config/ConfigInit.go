package config

import (
	"awesomeProxy/Log"
	"awesomeProxy/Report"
	"awesomeProxy/ac_automaton"
	"awesomeProxy/balance"
	"encoding/json"
	"io"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ProxyMethod bool    `json:"ProxyMethod"`
	ReProxy     ReProxy `json:"ReProxy"`
	CoProxy     CoProxy `json:"CoProxy"`
	Logg        Logg    `json:"Logg"`
}
type Backend struct {
	Host   string `json:"host"`
	Weight int    `json:"Weight"`
}
type Cache struct {
	Start   string `json:"start"`
	MaxSize string `json:"MaxSize"`
}
type ReProxy struct {
	Port          string    `json:"port"`
	BalanceMethod int       `json:"BalanceMethod"`
	Backend       []Backend `json:"backend"`
	Cache         Cache     `json:"Cache"`
}
type CoProxy struct {
	Port           string   `json:"port"`
	MultiListenNum int      `json:"MultiListenNum"`
	Nagle          bool     `json:"nagle"`
	Filt           []string `json:"filt"`
	WindowsUnProxy string   `json:"WindowsUnProxy"`
}
type SizeSplit struct {
	LogSize int    `json:"LogSize"`
	Unit    string `json:"Unit"`
	FileNum int    `json:"FileNum"`
}
type Logg struct {
	FileNameReProxy string    `json:"FileNameReProxy"`
	FileNameCoProxy string    `json:"FileNameCoProxy"`
	SplitFormat     string    `json:"SplitFormat"`
	DateSplit       string    `json:"DateSplit"`
	SizeSplit       SizeSplit `json:"SizeSplit"`
}

var CONFIG = &Config{}
var BalanceNames = []string{"hash", "random", "roundrobin", "weight_roundrobin", "shuffle", "shuffle2"}
var Insts []*balance.Instance
var Addrs []string

func (cc Config) Init() {
	ReadConfig()
	LogInit()
	Report.ReadReConfig()
	Report.ReadCoConfig()
	if CONFIG.ProxyMethod == true {
		cc.BalanceInit()
	} else {
		ac_automaton.AddKeyWordFromPath(ac_automaton.Acauto, "./ac_automaton/SensitiveWordsSmall.txt")
	}
	Logo()

}

// ReadConfig 读取配置文件
func ReadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		Log.Fatal("open json file err")
	}
	defer func() {
		if err := file.Close(); err != nil {
			Log.Error("关闭文件失败:", err)
		}
	}()
	data, err := io.ReadAll(file)
	if err != nil {
		Log.Fatal("read json file err")
	}
	err = json.Unmarshal(data, &CONFIG)
	if err != nil {
		Log.Fatal("Unmarshal json file err")
		return
	}
}

// BalanceInit 负载均衡初始化
func (cc Config) BalanceInit() {
	//这里获取负载均衡服务器ip
	BacSize := len(CONFIG.ReProxy.Backend)
	if BacSize == 0 {
		Log.Fatal("balance_init error,BacSize = 0")
	} else if BacSize == 1 {
		//如果server只有一个
		CONFIG.ReProxy.BalanceMethod = 2
	}
	for i := 0; i < BacSize; i++ {
		s := strings.Split(CONFIG.ReProxy.Backend[i].Host, ":")
		bhost := s[0]
		if len(s) != 2 {
			Log.Fatal("Config Backend Host not Equal 2")
		}
		bport, err := strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			Log.Fatal("balance_init trans int64 error")
		}
		wc := int64(CONFIG.ReProxy.Backend[i].Weight)
		one := balance.NewInstance(bhost, bport, wc)
		Insts = append(Insts, one)
	}
}

// LogInit 日志初始化
func LogInit() {
	LogName := ""
	// 设置日志目录
	DataPath, _ := os.Getwd()
	DataPath = DataPath + "\\Log\\Data"

	if CONFIG.ProxyMethod {
		LogName = CONFIG.Logg.FileNameReProxy
		DataPath = DataPath + "\\ReProxyData"
	} else {
		LogName = CONFIG.Logg.FileNameCoProxy
		DataPath = DataPath + "\\CoProxyData"
	}
	// 按照日期分割
	if CONFIG.Logg.SplitFormat == "DateSplit" {
		Method := Log.MODE_DAY
		if CONFIG.Logg.DateSplit == "MODE_DAY" {
			Method = Log.MODE_DAY
		} else if CONFIG.Logg.DateSplit == "MODE_HOUR" {
			Method = Log.MODE_HOUR
		} else if CONFIG.Logg.DateSplit == "MODE_MONTH" {
			Method = Log.MODE_MONTH
		}
		if CONFIG.ProxyMethod {

		}
		_, err := Log.SetRollingByTime(DataPath, LogName, Method)
		if err != nil {
			Log.Fatal("Log.SetRollingByTime error" + err.Error())
		}
		// 按照文件大小进行分割
	} else if CONFIG.Logg.SplitFormat == "SizeSplit" {
		Method := Log.MB
		if CONFIG.Logg.SizeSplit.Unit == "MB" {
			Method = Log.MB
		} else if CONFIG.Logg.SizeSplit.Unit == "KB" {
			Method = Log.KB
		} else if CONFIG.Logg.SizeSplit.Unit == "GB" {
			Method = Log.GB
		} else if CONFIG.Logg.SizeSplit.Unit == "TB" {
			Method = Log.TB
		}
		_, err := Log.SetRollingFileLoop(DataPath, LogName, int64(CONFIG.Logg.SizeSplit.LogSize), Method, CONFIG.Logg.SizeSplit.FileNum)
		if err != nil {
			Log.Fatal("Log.SetRollingByTime error" + err.Error())
		}
	} else {
		Log.Fatal("Log SplitFormat error")
	}
}

func Logo() {
	logo := `
'   ________   ________                   ________   ________   ________      ___    ___  ___    ___ 
'  |\   __  \ |\   ____\                 |\   __  \ |\   __  \ |\   __  \    |\  \  /  /||\  \  /  /|
'  \ \  \|\  \\ \  \___|_   ____________ \ \  \|\  \\ \  \|\  \\ \  \|\  \   \ \  \/  / /\ \  \/  / /
'   \ \   __  \\ \_____  \ |\____________\\ \   ____\\ \   _  _\\ \  \\\  \   \ \    / /  \ \    / / 
'    \ \  \ \  \\|____|\  \\|____________| \ \  \___| \ \  \\  \|\ \  \\\  \   /     \/    \/  /  /  
'     \ \__\ \__\ ____\_\  \                \ \__\     \ \__\\ _\ \ \_______\ /  /\   \  __/  / /    
'      \|__|\|__||\_________\                \|__|      \|__|\|__| \|_______|/__/ /\ __\|\___/ /     
'                \|_________|                                                |__|/ \|__|\|___|/      
'                                                                                                    
'                                                                                                     
`
	Log.Info(logo)
	Log.Info("欢迎使用awesomeProxy")
	if CONFIG.ProxyMethod {
		Log.Info("反向代理监听端口:0.0.0.0:" + CONFIG.ReProxy.Port)
	} else {
		Log.Info("正向代理监听端口:0.0.0.0:" + CONFIG.CoProxy.Port)
	}

}

// 分布式缓存组，下面是测试代码
// curl http://localhost:9090/ascache/scores/Tom
// curl http://localhost:9090/ascache/scores/kkk

// CacheInit 缓存初始化
func CacheInit() {
	for _, v := range CONFIG.ReProxy.Backend {
		uu := "http://" + v.Host
		Addrs = append(Addrs, uu)
	}
}
