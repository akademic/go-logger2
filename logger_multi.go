package logger

type MultiLogger struct {
	loggers []BaseLogger
}

func NewMultiLogger(loggers ...BaseLogger) BaseLogger {
	return &MultiLogger{loggers}
}

func (ml *MultiLogger) Print(v ...any) {
	for _, l := range ml.loggers {
		go l.Print(v...)
	}
}
