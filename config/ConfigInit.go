package config

import (
	"awesomeProxy/Cert"
	"awesomeProxy/Log"
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
}
type Backend struct {
	Host   string `json:"host"`
	Weghit int    `json:"Weghit"`
}
type ReProxy struct {
	Port          string    `json:"port"`
	BalanceMethod int       `json:"BalanceMethod"`
	Backend       []Backend `json:"backend"`
}
type CoProxy struct {
	Port           string   `json:"port"`
	MultiListenNum int      `json:"MultiListenNum"`
	Nagle          bool     `json:"nagle"`
	Filt           []string `json:"filt"`
}

var CONFIG = &Config{}
var BalanceNames = []string{"hash", "random", "roundrobin", "weight_roundrobin", "shuffle", "shuffle2"}
var Insts []*balance.Instance

func (cc Config) Init() {
	file, err := os.Open("config.json")
	if err != nil {
		Log.Log.Fatal("Config Init error : open json file err")
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		Log.Log.Fatal("Config Init error : read json file err")
		Log.Sfatal(err, "awesomeProxy/config/Init read file")
	}
	err = json.Unmarshal(data, &CONFIG)
	if err != nil {
		Log.Log.Fatal("Config Init error : Unmarshal json file err")
		return
	}
	if CONFIG.ProxyMethod == true {
		cc.balance_init()
	} else {
		Cert.Install()
	}

}

// 负载均衡初始化
func (cc Config) balance_init() {
	//这里获取负载均衡的分散ip
	BacSize := len(CONFIG.ReProxy.Backend)
	if BacSize == 0 {
		Log.Log.Fatal("awesomeProxy/config/Init balance_init error,BacSize = 0")
	} else if BacSize == 1 {
		//如果server只有一个
		CONFIG.ReProxy.BalanceMethod = 2
	}
	for i := 0; i < BacSize; i++ {
		s := strings.Split(CONFIG.ReProxy.Backend[i].Host, ":")
		bhost := s[0]
		if len(s) > 2 {
			Log.Log.Fatal("awesomeProxy/config/Init balance_init error,Backend Host error ")
		}
		if len(s) == 1 {
			s = append(s, "0")
		}

		bport, err := strconv.ParseInt(s[1], 10, 64)
		if err != nil {
			Log.Log.Fatal("balance_init trans int64 error")
		}

		wc := int64(CONFIG.ReProxy.Backend[i].Weghit)
		one := balance.NewInstance(bhost, bport, wc)
		Insts = append(Insts, one)
	}
}
