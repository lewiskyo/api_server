package number

import (
	"math/rand"
	"time"
)

//RandNum 获取 min~max之间的随机数
// @param min
// @param max
// @return int64
func RandNum(min, max int64) int64 {
	if min >= max {
		return max
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Int63n(max-min) + min // min~max
}

//GetMicroSecond
// @return int64
func GetMicroSecond() int64 {
	return time.Now().UnixNano() / 1e3
}

//GetMilliSecond
// @return int64
func GetMilliSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

//GetTimeStamp
// @return int64
func GetTimeStamp() int64 {
	return time.Now().Unix()
}
