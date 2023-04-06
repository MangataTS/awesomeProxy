package global

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type ReReport struct {
	RequestData   []RequestData `json:"RequestData"`
	CacheData     CacheData     `json:"CacheData"`
	BanIPReqTimes int64         `json:"BanIpReqTimes"`
	LogsData      LogsData      `json:"LogsData"`
	CrawlerData   []CrawlerData `json:"CrawlerData"`
	ServerStatus  ServerStatus  `json:"ServerStatus"`
}
type RequestData struct {
	Path  string `json:"path"`
	Times int64  `json:"times"`
}
type CacheData struct {
	ReqTimes  int64   `json:"ReqTimes"`
	HintTimes int64   `json:"HintTimes"`
	HintRate  float64 `json:"HintRate"`
	Memory    int64   `json:"memory"`
}
type LogsData struct {
	InfoTimes  int64 `json:"InfoTimes"`
	DebugTimes int64 `json:"DebugTimes"`
	WarnTimes  int64 `json:"WarnTimes"`
	ErrorTimes int64 `json:"ErrorTimes"`
	FatalTimes int64 `json:"FatalTimes"`
}
type CrawlerData struct {
	IP       string `json:"ip"`
	ReqTimes int64  `json:"ReqTimes"`
	BanTimes int64  `json:"BanTimes"`
}
type CPUINFO struct {
	Idle   float64 `json:"idle"`
	Kernel float64 `json:"kernel"`
	User   float64 `json:"user"`
	Usage  float64 `json:"Usage"`
}
type MEMORYINFO struct {
	Total     int64   `json:"total"`
	Available int64   `json:"available"`
	Usage     float64 `json:"usage"`
}
type DISKINFO struct {
	Total     int64   `json:"total"`
	Available int64   `json:"available"`
	Usage     float64 `json:"usage"`
}
type ServerStatus struct {
	CPUINFO    CPUINFO    `json:"CPU_INFO"`
	MEMORYINFO MEMORYINFO `json:"MEMORY_INFO"`
	DISKINFO   DISKINFO   `json:"DISK_INFO"`
}

// ReReportConfig 全局反向代理变量
var ReReportConfig = &ReReport{}

type CoReport struct {
	CoRequestData   []CoRequestData   `json:"CoRequestData"`
	CoProtocolData  []CoProtocolData  `json:"CoProtocolData"`
	CoBlackHostData []CoBlackHostData `json:"CoBlackHostData"`
	CoSensitiveData CoSensitiveData   `json:"CoSensitiveData"`
	LogsData        LogsData          `json:"LogsData"`
	ServerStatus    ServerStatus      `json:"ServerStatus"`
}
type CoRequestData struct {
	ReqHost  string `json:"reqHost"`
	ReqTimes int    `json:"reqTimes"`
}
type CoProtocolData struct {
	Name        string `json:"Name"`
	ReqTimes    int    `json:"ReqTimes"`
	ReqDataSize int    `json:"ReqDataSize"`
}
type CoBlackHostData struct {
	URLHost  string `json:"UrlHost"`
	ReqTimes int    `json:"ReqTimes"`
}
type CoSensitiveData struct {
	TriggerNum    int      `json:"TriggerNum"`
	Interceptions int      `json:"Interceptions"`
	IllegalURL    []string `json:"IllegalUrl"`
}

// Glock 全局锁
var Glock sync.Mutex

var CoReportConfig = &CoReport{}

var CalCoRequestData = make(map[string]int)
var CalCoProtocolData = make(map[string]CoProtocolData)
var CalCoBlackHostData = make(map[string]int)
var CalCoSensitiveDataUrl = make(map[string]bool)

func WriteToCalCoRequestData() {
	var WLock sync.Mutex
	WLock.Lock()
	defer WLock.Unlock()

	for host, cnt := range CalCoRequestData {
		flg := true
		for idx, item := range CoReportConfig.CoRequestData {
			if item.ReqHost == host {
				flg = false
				CoReportConfig.CoRequestData[idx].ReqTimes = cnt
				break
			}
		}
		if flg {
			CoReportConfig.CoRequestData = append(CoReportConfig.CoRequestData, CoRequestData{host, cnt})
		}
	}
}

func WriteCalCoProtocolData() {
	var WLock sync.Mutex
	WLock.Lock()
	defer WLock.Unlock()

	for idx, item := range CoReportConfig.CoProtocolData {
		value, ok := CalCoProtocolData[item.Name]
		if ok && value != CoReportConfig.CoProtocolData[idx] {
			CoReportConfig.CoProtocolData[idx] = value
		}
	}
}

func WriteCalCoBlackHostData() {
	var WLock sync.Mutex
	WLock.Lock()
	defer WLock.Unlock()

	for host, cnt := range CalCoBlackHostData {
		flg := true
		for idx, item := range CoReportConfig.CoBlackHostData {
			if item.URLHost == host {
				flg = false
				CoReportConfig.CoBlackHostData[idx].ReqTimes = cnt
				break
			}
		}
		if flg {
			CoReportConfig.CoBlackHostData = append(CoReportConfig.CoBlackHostData, CoBlackHostData{host, cnt})
		}
	}
}

func WriteCalCoSensitiveDataUrl() {
	var WLock sync.Mutex
	WLock.Lock()
	defer WLock.Unlock()

	for _, item := range CoReportConfig.CoSensitiveData.IllegalURL {
		value, ok := CalCoSensitiveDataUrl[item]
		if ok && value {
			CalCoSensitiveDataUrl[item] = false
		}
	}
	for name, cc := range CalCoSensitiveDataUrl {
		if cc {
			CoReportConfig.CoSensitiveData.IllegalURL = append(CoReportConfig.CoSensitiveData.IllegalURL, name)
		}
	}
}

func WriteCalLog(name string) {
	var WLock sync.Mutex
	WLock.Lock()
	defer WLock.Unlock()

	switch name {
	case "Debug":
		ReReportConfig.LogsData.DebugTimes++
		CoReportConfig.LogsData.DebugTimes++
	case "Info":
		ReReportConfig.LogsData.InfoTimes++
		CoReportConfig.LogsData.InfoTimes++
	case "Warn":
		ReReportConfig.LogsData.WarnTimes++
		CoReportConfig.LogsData.WarnTimes++
	case "Error":
		ReReportConfig.LogsData.ErrorTimes++
		CoReportConfig.LogsData.ErrorTimes++
	case "Fatal":
		ReReportConfig.LogsData.FatalTimes++
		CoReportConfig.LogsData.FatalTimes++
	}
	SaveReConfig()
	SaveCoConfig()

}

func SaveReConfig() {
	path := "./Report/Re/DataFile.json"
	data, err := json.MarshalIndent(ReReportConfig, "", " ")
	if err != nil {
		log.Println("json.MarshalIndent err ", err)
	}
	var Flock sync.Mutex
	Flock.Lock()
	defer Flock.Unlock()
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		log.Println("os.WriteFile ", err)
	}
}

func SaveCoConfig() {
	WriteToCalCoRequestData()
	WriteCalCoProtocolData()
	WriteCalCoBlackHostData()
	WriteCalCoSensitiveDataUrl()

	path := "./Report/Co/DataFile.json"
	data, err := json.MarshalIndent(CoReportConfig, "", " ")
	if err != nil {
		log.Println("json.MarshalIndent err ", err)
	}
	var Flock sync.Mutex
	Flock.Lock()
	defer Flock.Unlock()
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		log.Println("os.WriteFile ", err)
	}
}
