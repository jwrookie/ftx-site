package log

import (
	"path/filepath"

	"github.com/foxdex/ftx-site/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
)

func InitLog() {
	ws, level := getConfigLogArgs()
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConf)
	log := zap.New(
		zapcore.NewCore(encoder, ws, zap.NewAtomicLevelAt(level)),
		zap.AddCaller(),
	)
	Log = log
	Sugar = log.Sugar()
}

func getConfigLogArgs() (zapcore.WriteSyncer, zapcore.Level) {
	log := config.GetConfig().Log
	level := zap.InfoLevel

	switch log.Level {
	case "ERROR", "error":
		level = zap.ErrorLevel
	case "WARN", "warn":
		level = zap.WarnLevel
	case "", "INFO", "info":
		level = zap.InfoLevel
	}

	var syncers []zapcore.WriteSyncer
	if log.FileName != "" {
		logger := &lumberjack.Logger{
			// if logs dir not exist, it will be auto create
			Filename:   filepath.Join(log.Dir, log.FileName),
			MaxSize:    log.MaxSize,
			MaxBackups: log.MaxBackups,
			MaxAge:     log.MaxAge,
			Compress:   log.Compress,
			LocalTime:  true,
		}
		syncers = append(syncers, zapcore.AddSync(logger))
	}

	ws := zapcore.NewMultiWriteSyncer(syncers...)

	return ws, level
}
