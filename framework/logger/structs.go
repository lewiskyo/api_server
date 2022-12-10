package logger

import (
	"context"
	"io"

	"go.uber.org/zap"
)

const CtxKeyTraceId = "tid"

type Logger struct {
	*zap.Logger
	w io.Writer
}

//NewWithZap
// @param z
// @return *Logger
func NewWithZap(z *zap.Logger, writer io.Writer) *Logger {
	return &Logger{Logger: z, w: writer}
}

//Zap
// @receiver log
// @return *zap.Logger
func (log *Logger) clone() *Logger {
	c := *log
	return &c
}

// WithCtx
// return new logger with context
func (log *Logger) WithCtx(ctx context.Context) (l *Logger) {
	if ctx == nil {
		return log
	}
	l = log.clone()

	var fields []zap.Field
	if traceId, ok := extractTraceID(ctx); ok {
		fields = append(fields, zap.String(CtxKeyTraceId, traceId))
	}

	if len(fields) > 0 {
		l.Logger = l.Logger.With(fields...)
	}

	return l
}

//Zap
// @receiver log
// @return *zap.Logger
func (log *Logger) Zap() *zap.Logger {
	return log.Logger
}

//Writer
// @receiver log
// @return io.Writer
func (log *Logger) Writer() io.Writer {
	return log.w
}

//Debugf
// @receiver log
// @param format
// @param args
func (log *Logger) Debugf(format string, args ...interface{}) {
	log.Sugar().Debugf(format, args...)
}

//Infof
// @receiver log
// @param format
// @param args
func (log *Logger) Infof(format string, args ...interface{}) {
	log.Sugar().Infof(format, args...)
}

//Warnf
// @receiver log
// @param format
// @param args
func (log *Logger) Warnf(format string, args ...interface{}) {
	log.Sugar().Warnf(format, args...)
}

//Panicf
// @receiver log
// @param format
// @param args
func (log *Logger) Panicf(format string, args ...interface{}) {
	log.Sugar().Panicf(format, args...)
}

//Panicf
// @receiver log
// @param format
// @param args
func (log *Logger) DPanicf(format string, args ...interface{}) {
	log.Sugar().DPanicf(format, args...)
}

//Errorf
// @receiver log
// @param format
// @param args
func (log *Logger) Errorf(format string, args ...interface{}) {
	log.Sugar().Errorf(format, args...)
}

func (log *Logger) sugar() *zap.Logger {
	return log.Logger
}

func (log *Logger) Debug(msg string, fields ...zap.Field) {
	log.sugar().Debug(msg, fields...)
}

func (log *Logger) Info(msg string, fields ...zap.Field) {
	log.sugar().Info(msg, fields...)
}

func (log *Logger) Warn(msg string, fields ...zap.Field) {
	log.sugar().Warn(msg, fields...)
}

func (log *Logger) Error(msg string, fields ...zap.Field) {
	log.sugar().Error(msg, fields...)
}

func (log *Logger) DPanic(msg string, fields ...zap.Field) {
	log.sugar().DPanic(msg, fields...)
}

func (log *Logger) Panic(msg string, fields ...zap.Field) {
	log.sugar().Panic(msg, fields...)
}

func (log *Logger) Fatal(msg string, fields ...zap.Field) {
	log.sugar().Fatal(msg, fields...)
}

func (log *Logger) With(fields ...zap.Field) (l *Logger) {
	return log.WithZapFields(fields...)
}

func (log *Logger) WithZapFields(fields ...zap.Field) (l *Logger) {
	l = log.clone()

	if len(fields) > 0 {
		l.Logger = l.Logger.With(fields...)
	}

	return l
}
