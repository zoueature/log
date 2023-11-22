package log

import (
	"context"
	"github.com/google/uuid"
	"github.com/jiebutech/config"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

const logIdKey = "log-id"

func Init(cfg *config.LogConfig) error {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	writer, err := cfg.GetLogWriter()
	if err != nil {
		return err
	}
	logrus.SetOutput(writer)
	logger = logrus.StandardLogger()
	return nil
}

type setter interface {
	Set(key string, value interface{})
}

// InjectLogID 注入logid到上下文中
func InjectLogID(ctx setter) {
	logIdVal, _ := uuid.NewUUID()
	ctx.Set(logIdKey, logIdVal)
}

// Debug 输出Debug日志
func Debug(ctx context.Context, args ...interface{}) {
	logger.WithField(logIdKey, ctx.Value(logIdKey)).Debug(args...)
}

// Debugf 输出Debug日志
func Debugf(ctx context.Context, format string, args ...interface{}) {
	logger.WithField(logIdKey, ctx.Value(logIdKey)).Debugf(format, args...)
}

// Info 输出Info日志
func Info(ctx context.Context, args ...interface{}) {
	logger.WithField(logIdKey, ctx.Value(logIdKey)).Info(args...)
}

// Infof 输出Info日志
func Infof(ctx context.Context, format string, args ...interface{}) {
	logger.WithField(logIdKey, ctx.Value(logIdKey)).Infof(format, args...)
}

// Warn 输出Warn日志
func Warn(ctx context.Context, args ...interface{}) {
	logger.WithField(logIdKey, ctx.Value(logIdKey)).Warn(args...)
}

// Warnf 输出Warn日志
func Warnf(ctx context.Context, format string, args ...interface{}) {
	logger.WithField(logIdKey, ctx.Value(logIdKey)).Warnf(format, args...)
}

// Error 输出error日志
func Error(ctx context.Context, args ...interface{}) {
	logger.WithField(logIdKey, ctx.Value(logIdKey)).Error(args...)
}

// Errorf 输出error日志
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logger.WithField(logIdKey, ctx.Value(logIdKey)).Errorf(format, args...)
}

// Panic 输出Panic日志
func Panic(ctx context.Context, args ...interface{}) {
	logger.WithField(logIdKey, ctx.Value(logIdKey)).Panic(args...)
}

// Panicf 输出Panic日志
func Panicf(ctx context.Context, format string, args ...interface{}) {
	logger.WithField(logIdKey, ctx.Value(logIdKey)).Panicf(format, args...)
}
