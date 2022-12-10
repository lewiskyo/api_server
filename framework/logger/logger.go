package logger

import (
	"api_server/framework/config"
	"context"
	"os"
	"strings"
	"sync"
	"time"
	"unsafe"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// errorLogger
	errorLogger unsafe.Pointer //error_log
	// accessLogger
	accessLogger unsafe.Pointer //access_log
	// serverLogger 请使用 GetServerLogger 方法
	serverLogger unsafe.Pointer //server_log

	//timezoneSyncOnce
	timezoneSyncOnce sync.Once
	//timezone time zone location
	timezone *time.Location
)

//获取日志级别 #debug:-1 info:0 warn:1 error:2 dpanic:3 panic:4 fatal:5
//默认为 INFO
func getLogLevel(level string) zapcore.Level {
	level = strings.ToLower(level)
	levelInt := zapcore.InfoLevel
	if level == "debug" {
		levelInt = zapcore.DebugLevel
	} else if level == "info" {
		levelInt = zapcore.InfoLevel
	} else if level == "warn" {
		levelInt = zapcore.WarnLevel
	} else if level == "error" {
		levelInt = zapcore.ErrorLevel
	} else if level == "dpanic" {
		levelInt = zapcore.DPanicLevel
	} else if level == "panic" {
		levelInt = zapcore.PanicLevel
	} else if level == "fatal" {
		levelInt = zapcore.FatalLevel
	}
	return levelInt
}

//日志文件初始化
func init() {
	var l *Logger
	conf := config.LogConfig

	// access log
	l = getZapLogger(conf.AccessLog, conf, "info", false, config.DisableAccessLogRecord(), conf.OutputStdout&1 == 1)
	accessLogger = unsafe.Pointer(l)

	// error log
	l = getZapLogger(conf.ErrorLog, conf, conf.ErrorLogLevel, conf.ErrorLogFullCaller, false, conf.OutputStdout&2 == 2)
	errorLogger = unsafe.Pointer(l)

	// server log
	l = getZapLogger(conf.ServerLog, conf, conf.ServerLogLevel, false, false, conf.OutputStdout&4 == 4)
	serverLogger = unsafe.Pointer(l)

}

// FlushLogger flush log to disk
func FlushLogger() {
	flushable := []*Logger{
		Access(),
		Err(),
		Server(),
	}

	for _, l := range flushable {
		if l == nil {
			continue
		}
		_ = l.Sync()
	}
}

var defaultZapOpts = []zap.Option{
	zap.AddCaller(),
	zap.AddCallerSkip(1),
}

//resolveTimeZone parse timezone
// @return *time.Location
func resolveTimeZone() *time.Location {
	timezoneSyncOnce.Do(func() {
		if len(config.LogConfig.TimeZone) <= 0 {
			timezone = time.Local
			return
		}
		// Get timezone location
		var err error
		timezone, err = time.LoadLocation(config.LogConfig.TimeZone)
		if err != nil || timezone == nil {
			timezone = time.Local
		}
	})
	return timezone
}

//getZapLogger
// @param filePath
// @param conf
// @param logLevel
// @param fullCallerEncoder
// @param disable
// @param outputStdout
// @return *Logger
func getZapLogger(filePath string, conf *config.LogConf, logLevel string, fullCallerEncoder, disable, outputStdout bool) *Logger {
	var syncer zapcore.WriteSyncer
	if disable {
		// output to a black hole
		syncer = getBlackHole()
	} else {
		syncer = getLogWriter(filePath, conf.MaxSize, conf.MaxBackups, conf.MaxAge, conf.Compress)
		if outputStdout {
			syncer = zap.CombineWriteSyncers(syncer, os.Stdout)
		}
	}

	encoder := getEncoder(fullCallerEncoder)

	level := getLogLevel(logLevel)
	if level < zapcore.DebugLevel || level > zapcore.FatalLevel {
		level = zapcore.InfoLevel
	}

	z := zap.New(
		zapcore.NewCore(
			encoder,
			syncer,
			level,
		),
		defaultZapOpts...,
	)
	return NewWithZap(z, syncer)
}

// getEncoder 获取日志编码格式
// @return zapcore.Encoder
func getEncoder(fullCallerEncoder bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	//日志时间格式
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.In(resolveTimeZone()).Format(config.LogConfig.TimeFormat))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if fullCallerEncoder {
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

//获取日志写入器
// getLogWriter
// @param filePath
// @param maxSize
// @param maxBackups
// @param maxAge
// @param compress
// @return zapcore.WriteSyncer
func getLogWriter(filePath string, maxSize, maxBackups, maxAge int, compress bool) zapcore.WriteSyncer {
	hook := &lumberjack.Logger{
		Filename:   filePath,   //filePath
		MaxSize:    maxSize,    // megabytes
		MaxBackups: maxBackups, //backups
		MaxAge:     maxAge,     //days
		Compress:   compress,   // disabled by default
	}
	return zapcore.AddSync(hook)
}

func getBlackHole() zapcore.WriteSyncer {
	hook := &BlackHole{}
	return zapcore.AddSync(hook)
}

// Server
// server_log 日志 logger
func Server() *Logger {
	return (*Logger)(serverLogger)
}

// GetServerLogger
// server_log 日志 logger
func GetServerLogger() *Logger {
	return Server()
}

// Err
// error_log 日志 logger
func Err() (l *Logger) {
	return (*Logger)(errorLogger)
}

// GetErrorLogger
// error_log 日志 logger
func GetErrorLogger() *Logger {
	return Err()
}

// Access
// access_log 日志 logger
func Access() *Logger {
	return (*Logger)(accessLogger)
}

// GetAccessLogger
// access_log 日志 logger
func GetAccessLogger() *Logger {
	return Access()
}

// ##############################
// errorLog method shortcuts

func Debug(args ...interface{}) {
	Err().Sugar().Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	Err().Sugar().Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	Err().Sugar().Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	Err().Sugar().Info(args...)
}

func Infof(template string, args ...interface{}) {
	Err().Sugar().Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	Err().Sugar().Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	Err().Sugar().Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	Err().Sugar().Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	Err().Sugar().Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	Err().Sugar().Error(args...)
}

func Errorf(template string, args ...interface{}) {
	Err().Sugar().Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	Err().Sugar().Errorw(msg, keysAndValues...)
}

func DPanic(args ...interface{}) {
	Err().Sugar().DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	Err().Sugar().DPanicf(template, args...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	Err().Sugar().DPanicw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	Err().Sugar().Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	Err().Sugar().Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	Err().Sugar().Panicw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	Err().Sugar().Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	Err().Sugar().Fatalf(template, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	Err().Sugar().Fatalw(msg, keysAndValues...)
}

func WithCtx(ctx context.Context) (l *Logger) {
	return Err().WithCtx(ctx)
}

// extractTraceID HTTP使用request.Context，不要使用错了
func extractTraceID(ctx context.Context) (string, bool) {
	span := trace.SpanContextFromContext(ctx)
	if span.HasTraceID() {
		return span.TraceID().String(), true
	}
	return "", false
}
