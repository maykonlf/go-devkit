package log

import (
	"go.uber.org/zap"
)

type LoggerI interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	DPanic(msg string, keysAndValues ...interface{})
	Panic(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
}

func NewLogger(level Level) LoggerI {
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(levelMap[level]),
		Encoding:         "json",
		EncoderConfig:    EncoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := config.Build()
	return &Logger{
		logger: logger.Sugar(),
	}
}

type Logger struct {
	logger *zap.SugaredLogger
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *Logger) DPanic(msg string, keysAndValues ...interface{}) {
	l.logger.DPanicw(msg, keysAndValues...)
}

func (l *Logger) Panic(msg string, keysAndValues ...interface{}) {
	l.logger.Panicw(msg, keysAndValues...)
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}
