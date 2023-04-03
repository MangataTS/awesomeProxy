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
	file, err := os.Open("./Report/Re/DataFile.json")
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

// ReadCoConfig 读取正向代理报告配置文件
func ReadCoConfig() {
	file, err := os.Open("./Report/Co/DataFile.json")
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
	err = json.Unmarshal(data, &global.CoReportConfig)
	if err != nil {
		Log.Fatal("Unmarshal json file err: ", err)
		return
	}
}
