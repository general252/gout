package ulog2

type LogLevel int

// RFC5424 log message levels.
const (
	LevelEmergency LogLevel = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

const (
	LevelInfo  = LevelInformational
	LevelTrace = LevelDebug
	LevelWarn  = LevelWarning
)

func (v LogLevel) String() string {
	switch v {
	case LevelEmergency:
		return "Emergency"
	case LevelAlert:
		return "Alert"
	case LevelCritical:
		return "Critical"
	case LevelError:
		return "E"
	case LevelWarning:
		return "W"
	case LevelNotice:
		return "Notice"
	case LevelInformational:
		return "I"
	case LevelDebug:
		return "D"
	default:
		return "UnKnown"
	}
}

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
