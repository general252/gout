package ulog2

func Debug(v ...interface{}) {
	Component().Debug(v...)
}

func Info(v ...interface{}) {
	Component().Info(v...)
}

func Warn(v ...interface{}) {
	Component().Warn(v...)
}

func Error(v ...interface{}) {
	Component().Error(v...)
}

func DebugT(tags Tag, v ...interface{}) {
	Component(tags...).Debug(v...)
}

func InfoT(tags Tag, v ...interface{}) {
	Component(tags...).Info(v...)
}

func WarnT(tags Tag, v ...interface{}) {
	Component(tags...).Warn(v...)
}

func ErrorT(tags Tag, v ...interface{}) {
	Component(tags...).Error(v...)
}
