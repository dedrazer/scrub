package utils

import "go.uber.org/zap/zapcore"

const (
	LevelDebug = "debug"
	LevelErr   = "err"
	LevelInfo  = "info"
	LevelWarn  = "warn"
)

func ParseZapLevel(level string) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelErr:
		return zapcore.ErrorLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	default:
		return zapcore.InfoLevel
	}
}
