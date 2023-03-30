package global

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
	ReqTimes  int64 `json:"ReqTimes"`
	HintTimes int64 `json:"HintTimes"`
	HintRate  int64 `json:"HintRate"`
	Memory    int64 `json:"memory"`
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

// 存放全局变量
var ReReportConfig = &ReReport{}
