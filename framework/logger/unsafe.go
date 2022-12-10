package logger

import (
	"api_server/framework/config"
	"sync/atomic"
	"unsafe"
)

func UnsafeChangeErrorLogPath(newFilePath string) (rollback func()) {
	conf := config.LogConfig
	l := getZapLogger(newFilePath, conf, conf.ErrorLogLevel, conf.ErrorLogFullCaller, false, conf.OutputStdout&2 == 2)
	newLogger := unsafe.Pointer(l)
	old := atomic.SwapPointer(&errorLogger, newLogger)
	rollback = func() {
		atomic.SwapPointer(&errorLogger, old)
	}
	return
}
