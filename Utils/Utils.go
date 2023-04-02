package Utils

import (
	"awesomeProxy/Log"
	"awesomeProxy/config"
	"bytes"
	"crypto/tls"
	"net/http"
	"os"
	"reflect"
	"regexp"
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

func BlacklistFilter(request *http.Request) bool {
	Turl := request.URL.String()
	for _, s := range config.CONFIG.CoProxy.Filt {
		//Log.Debug(s)
		r, err := regexp.Compile(s)
		if err != nil {
			Log.Error("正则表达式有问题，请仔细检查，这里只支持Golang的正则~")
		}
		if r.MatchString(Turl) {
			request.URL.Host = "-1"
			request.Host = "-1"
			Log.Warn("访问到黑名单网页: " + s)
			return true
		}
	}
	return false
}
