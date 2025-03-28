package logger

import "fmt"

type Log struct {
	Level     LogLevel
	Component string
	Message   string
}

func (l Log) String() string {
	if l.Component == "" {
		return fmt.Sprintf("[%s]: %s", l.Level, l.Message)
	}

	return fmt.Sprintf("[%s] [%s]: %s", l.Level, l.Component, l.Message)
}

func (l Log) Labels() map[string]string {
	if l.Component == "" {
		return nil
	}

	return map[string]string{
		"component": l.Component,
	}
}
