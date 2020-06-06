package log

var (
	logger LoggerI
)

func init() {
	logger = NewLogger(JsonFormat, InfoLevel)
}

func Config(format Format, level Level) {
	logger = NewLogger(format, level)
}

func SetLogger(l LoggerI) {
	logger = l
}

func Debug(msg string, keysAndValues ...interface{}) {
	logger.Debug(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	logger.Info(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	logger.Warn(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	logger.Error(msg, keysAndValues...)
}

func DPanic(msg string, keysAndValues ...interface{}) {
	logger.DPanic(msg, keysAndValues...)
}

func Panic(msg string, keysAndValues ...interface{}) {
	logger.Panic(msg, keysAndValues...)
}
