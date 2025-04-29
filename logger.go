package logger

import (
	"fmt"
)

type BaseLogger interface {
	Print(v ...any)
}

type LoggerImpl struct {
	config     *Config
	component  string
	baseLogger BaseLogger
}

func New(baseLogger BaseLogger, component string, config *Config) *LoggerImpl {
	return &LoggerImpl{
		config:     config,
		component:  component,
		baseLogger: baseLogger,
	}
}

func (l *LoggerImpl) WithComponent(component string) Logger {
	return &LoggerImpl{
		config:     l.config,
		component:  component,
		baseLogger: l.baseLogger,
	}
}

func (l *LoggerImpl) Info(pattern string, args ...any) {
	if !l.logOn(LogInfo) {
		return
	}

	message := pattern
	if len(args) > 0 {
		message = fmt.Sprintf(pattern, args...)
	}

	log := Log{
		Level:     LogInfo,
		Component: l.component,
		Message:   message,
	}

	l.baseLogger.Print(log)
}

func (l *LoggerImpl) Debug(pattern string, args ...any) {
	if !l.logOn(LogDebug) {
		return
	}

	message := pattern
	if len(args) > 0 {
		message = fmt.Sprintf(pattern, args...)
	}

	log := Log{
		Level:     LogDebug,
		Component: l.component,
		Message:   message,
	}

	l.baseLogger.Print(log)
}

func (l *LoggerImpl) Error(pattern string, args ...any) {
	if !l.logOn(LogError) {
		return
	}

	message := pattern
	if len(args) > 0 {
		message = fmt.Sprintf(pattern, args...)
	}

	log := Log{
		Level:     LogError,
		Component: l.component,
		Message:   message,
	}

	l.baseLogger.Print(log)
}

func (l *LoggerImpl) SetConfig(config *Config) {
	if config == nil {
		return
	}

	l.config.Level = config.Level

	if config.ComponentLevel != nil {
		l.config.ComponentLevel = config.ComponentLevel
	}
}

func (l *LoggerImpl) logOn(level LogLevel) bool {
	if l.component == "" {
		return l.config.Level.CanLog(level)
	}

	confLevel, ok := l.config.ComponentLevel[l.component]
	if !ok {
		return l.config.Level.CanLog(level)
	}

	return confLevel.CanLog(level)
}
