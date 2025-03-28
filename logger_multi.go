package logger

import "sync"

type MultiLogger struct {
	loggers []BaseLogger
}

func NewMultiLogger(loggers ...BaseLogger) BaseLogger {
	return &MultiLogger{loggers}
}

func (ml *MultiLogger) Print(v ...any) {
	var wg sync.WaitGroup
	for _, l := range ml.loggers {
		wg.Add(1)
		go func() {
			l.Print(v...)
			wg.Done()
		}()
	}

	wg.Wait()
}
