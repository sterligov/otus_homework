package logger

type Logger interface {
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Named(string) Logger
}

var globalLogger Logger

func SetGlobalLogger(l Logger) {
	globalLogger = l
}

func Named(name string) Logger {
	return globalLogger.Named(name)
}

func Infof(format string, vals ...interface{}) {
	globalLogger.Infof(format, vals...)
}

func Debugf(format string, vals ...interface{}) {
	globalLogger.Debugf(format, vals...)
}

func Warnf(format string, vals ...interface{}) {
	globalLogger.Warnf(format, vals...)
}

func Errorf(format string, vals ...interface{}) {
	globalLogger.Errorf(format, vals...)
}
