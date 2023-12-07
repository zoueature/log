package log

import (
	"context"
	"github.com/google/uuid"
	"github.com/zoueature/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

var zapLogger *zap.Logger

const (
	logIdKey  = "log-id"
	userIdKey = "user-id"
	uriKey    = "uriKey"

	defaultLogPath = "./logs"
)

func init() {
	var err error
	zapLogger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

type apiCtx interface {
	AuthUserID() int
	RequestURI() string
}

type debugLevelEnabler string

func (e debugLevelEnabler) Enabled(level zapcore.Level) bool {
	return level <= zapcore.DebugLevel
}

type infoLevelEnabler string

func (e infoLevelEnabler) Enabled(level zapcore.Level) bool {
	return level > zapcore.DebugLevel && level < zapcore.ErrorLevel
}

type errorLevelEnabler string

func (e errorLevelEnabler) Enabled(level zapcore.Level) bool {
	return level >= zapcore.ErrorLevel
}

// Configure 配置log
func Configure(cfg *config.Configuration) error {
	if cfg.App.Debug {
		var err error
		zapLogger, err = zap.NewDevelopment()
		if err != nil {
			return err
		}
	} else {
		debugWs, err := getLogWriter(*cfg.Log, "debug")
		if err != nil {
			return err
		}
		infoWs, err := getLogWriter(*cfg.Log, "info")
		if err != nil {
			return err
		}
		errorWs, err := getLogWriter(*cfg.Log, "error")
		if err != nil {
			return err
		}
		prodconf := zap.NewProductionEncoderConfig()
		// 实例化日志器
		debugCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodconf), debugWs, debugLevelEnabler(""))
		infoCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodconf), infoWs, infoLevelEnabler(""))
		errCore := zapcore.NewCore(zapcore.NewJSONEncoder(prodconf), errorWs, errorLevelEnabler(""))

		zapLogger = zap.New(zapcore.NewTee(debugCore, infoCore, errCore))
		if err != nil {
			return err
		}
	}
	return nil
}

type setter interface {
	Set(key string, value interface{})
}

// getLogWriter 获取日志输出方式  日志文件 控制台
func getLogWriter(conf config.LogConfig, level string) (zapcore.WriteSyncer, error) {
	if conf.Driver == config.StdoutDriver {
		// 日志同时输出到控制台
		return zapcore.AddSync(os.Stdout), nil
	}

	// 判断日志路径是否存在，如果不存在就创建
	if exist := IsExist(conf.Path); !exist {
		if conf.Path == "" {
			conf.Path = defaultLogPath
		}
		if err := os.MkdirAll(conf.Path, os.ModePerm); err != nil {
			// 指定目录失败， 则写在默认目录
			conf.Path = defaultLogPath
			if err := os.MkdirAll(conf.Path, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}

	// 日志文件 与 日志切割 配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filepath.Join(conf.Path, level+".log"), // 日志文件路径
		MaxSize:    conf.LogFileMaxSize,                    // 单个日志文件最大多少 mb
		MaxBackups: conf.LogFileMaxBackups,                 // 日志备份数量
		MaxAge:     conf.LogMaxAge,                         // 日志最长保留时间
		Compress:   true,                                   // 是否压缩日志
	}
	// 日志只输出到日志文件
	return zapcore.AddSync(lumberJackLogger), nil
}

// IsExist 判断文件或者目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// InjectLogID 注入logid到上下文中
func InjectLogID(ctx setter) {
	logIdVal, _ := uuid.NewUUID()
	ctx.Set(logIdKey, logIdVal)
}

func injectCommonValue(ctx context.Context) *zap.Logger {
	ac, ok := ctx.(apiCtx)

	logger := zapLogger.With(zap.Any(logIdKey, ctx.Value(logIdKey)))
	if ok {
		logger = logger.With(zap.Any(userIdKey, ac.AuthUserID()), zap.Any(uriKey, ac.RequestURI()))
	}
	return logger
}

// Debug 输出Debug日志
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	injectCommonValue(ctx).Debug(msg, fields...)
}

// Info 输出Info日志
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	injectCommonValue(ctx).Info(msg, fields...)
}

// Warn 输出Warn日志
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	injectCommonValue(ctx).Warn(msg, fields...)
}

// Error 输出error日志
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	injectCommonValue(ctx).Error(msg, fields...)
}

// Panic 输出Panic日志
func Panic(ctx context.Context, msg string, fields ...zap.Field) {
	injectCommonValue(ctx).Panic(msg, fields...)
}
