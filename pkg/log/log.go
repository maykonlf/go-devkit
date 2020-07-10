package log

var (
	logger  LoggerI
	funcMap map[Level]func(msg string, keysAndValues ...interface{})
)

//nolint:gochecknoinits // proposital init() function to make package usage easier
func init() {
	logger = NewLogger(JSONFormat, InfoLevel)
	funcMap = map[Level]func(msg string, keysAndValues ...interface{}){
		DebugLevel:  Debug,
		InfoLevel:   Info,
		WarnLevel:   Warn,
		ErrorLevel:  Error,
		DPanicLevel: DPanic,
		PanicLevel:  Panic,
	}
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

func Log(level Level, msg string, keysAndValues ...interface{}) {
	funcMap[level](msg, keysAndValues...)
}
