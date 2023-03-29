package Utils

import (
	"awesomeProxy/Log"
	"awesomeProxy/config"
	"bytes"
	"crypto/tls"
	"fmt"
	"golang.org/x/sys/windows"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"syscall"
	"unsafe"
)

func FileExist(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func GetLastTimeFrame(conn *tls.Conn, property string) []byte {
	rawInputPtr := reflect.ValueOf(conn).Elem().FieldByName(property)
	if rawInputPtr.Kind() != reflect.Struct {
		return []byte{}
	}
	val, _ := reflect.NewAt(rawInputPtr.Type(), unsafe.Pointer(rawInputPtr.UnsafeAddr())).Elem().Interface().(bytes.Buffer)
	return val.Bytes()
}

func BlacklistFilter(request *http.Request) {
	Turl := request.URL.String()
	for _, s := range config.CONFIG.CoProxy.Filt {
		Log.Info(s)
		r, err := regexp.Compile(s)
		if err != nil {
			Log.Fatal("正则表达式有问题，请仔细检查，这里只支持Golang的正则~")
		}
		if r.MatchString(Turl) {
			request.URL.Host = "-1"
			request.Host = "-1"
			Log.Warn("访问到黑名单网页: " + s)
			return
		}
	}
}

const (
	INTERNET_PER_CONN_FLAGS               = 1
	INTERNET_PER_CONN_PROXY_SERVER        = 2
	INTERNET_PER_CONN_PROXY_BYPASS        = 3
	INTERNET_OPTION_REFRESH               = 37
	INTERNET_OPTION_SETTINGS_CHANGED      = 39
	INTERNET_OPTION_PER_CONNECTION_OPTION = 75
)

type INTERNET_PER_CONN_OPTION struct {
	dwOption uint32
	dwValue  uint64 // 注意 32位 和 64位 struct 和 union 内存对齐
}

type INTERNET_PER_CONN_OPTION_LIST struct {
	dwSize        uint32
	pszConnection *uint16
	dwOptionCount uint32
	dwOptionError uint32
	pOptions      uintptr
}

func SetProxy(proxy string) error {
	winInet, err := windows.LoadLibrary("Wininet.dll")
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("LoadLibrary Wininet.dll Error: %s", err))
	}
	InternetSetOptionW, err := windows.GetProcAddress(winInet, "InternetSetOptionW")
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("GetProcAddress InternetQueryOptionW Error: %s", err))
	}

	options := [3]INTERNET_PER_CONN_OPTION{}
	options[0].dwOption = INTERNET_PER_CONN_FLAGS
	if proxy == "" {
		options[0].dwValue = 1
	} else {
		options[0].dwValue = 2
	}
	options[1].dwOption = INTERNET_PER_CONN_PROXY_SERVER
	options[1].dwValue = uint64(uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(proxy))))
	options[2].dwOption = INTERNET_PER_CONN_PROXY_BYPASS
	options[2].dwValue = uint64(uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(config.CONFIG.CoProxy.WindowsUnProxy))))

	list := INTERNET_PER_CONN_OPTION_LIST{}
	list.dwSize = uint32(unsafe.Sizeof(list))
	list.pszConnection = nil
	list.dwOptionCount = 3
	list.dwOptionError = 0
	list.pOptions = uintptr(unsafe.Pointer(&options))

	// https://www.cnpython.com/qa/361707
	callInternetOptionW := func(dwOption uintptr, lpBuffer uintptr, dwBufferLength uintptr) error {
		r1, _, err := syscall.Syscall6(InternetSetOptionW, 4, 0, dwOption, lpBuffer, dwBufferLength, 0, 0)
		if r1 != 1 {
			return err
		}
		return nil
	}

	err = callInternetOptionW(INTERNET_OPTION_PER_CONNECTION_OPTION, uintptr(unsafe.Pointer(&list)), uintptr(unsafe.Sizeof(list)))
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_PER_CONNECTION_OPTION Error: %s", err)
	}
	err = callInternetOptionW(INTERNET_OPTION_SETTINGS_CHANGED, 0, 0)
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_SETTINGS_CHANGED Error: %s", err)
	}
	err = callInternetOptionW(INTERNET_OPTION_REFRESH, 0, 0)
	if err != nil {
		return fmt.Errorf("INTERNET_OPTION_REFRESH Error: %s", err)
	}
	return nil
}

// SetWindowsProxy 设置系统代理 eg: 127.0.0.1:9090
func SetWindowsProxy(Host string) {
	err := SetProxy(Host)
	if err != nil {
		Log.Fatal("Windows系统代理：", Host, "自动设置失败", err)
	}
	Log.Info("Windows系统代理：", Host, "自动设置成功")
}
