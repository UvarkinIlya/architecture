package logger

var (
	loggerInstance Logger
)

func ConfigurateLogger(logFile string) {
	loggerInstance, _ = New(logFile)
}

func Info(msg string, args ...interface{}) {
	loggerInstance.Info(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	loggerInstance.Debug(msg, args...)
}

func Error(msg string, args ...interface{}) {
	loggerInstance.Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	loggerInstance.Fatal(msg, args...)
}
