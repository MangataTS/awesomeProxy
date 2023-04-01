package Report

import (
	"awesomeProxy/Log"
	"awesomeProxy/global"
	"encoding/json"
	"io"
	"os"
)

// ReadReConfig 读取反向代理报告配置文件
func ReadReConfig() {
	path, err := os.Getwd()
	if err != nil {
		Log.Error("Os Getwd err")
	}
	path = path + "\\Report\\Re\\DataFile.json"
	file, err := os.Open(path)
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
		Log.Fatal("read json file err: ", err)
	}
	err = json.Unmarshal(data, &global.ReReportConfig)
	if err != nil {
		Log.Fatal("Unmarshal json file err: ", err)
		return
	}
}
