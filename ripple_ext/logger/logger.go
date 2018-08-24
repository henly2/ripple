package logger

type Logger interface {
	Info(arg0 interface{}, args ...interface{})
	Trace(arg0 interface{}, args ...interface{})
	Warn(arg0 interface{}, args ...interface{}) error
	Error(arg0 interface{}, args ...interface{}) error
	Debug(arg0 interface{}, args ...interface{})
}
