package strhelper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

//Md5 获取Md5值
// @param data
// @return string
func Md5(data interface{}) string {
	h := md5.New()
	stringData := fmt.Sprint(data)
	h.Write([]byte(stringData))
	return hex.EncodeToString(h.Sum(nil))
}

//GetCurrentDate 获取当前时间 YYYY-MM-DD HH:ii:SS
// @param format
// @return string
func GetCurrentDate(format ...string) string {
	var formatStr string
	if len(format) <= 0 {
		formatStr = "2006-01-02 15:04:05"
	} else {
		formatStr = format[0]
	}
	return time.Now().Format(formatStr)
}

//Empty 判断字符串是否为空
// @param str
// @return bool
func Empty(str string) bool {
	if len(strings.TrimSpace(str)) <= 0 {
		return true
	}
	return false
}
