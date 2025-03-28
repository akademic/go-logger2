// Package logger contains project-wide logging interface
package logger

type Logger interface {
	WithComponent(component string) Logger
	Info(pattern string, args ...any)
	Error(pattern string, args ...any)
	Debug(pattern string, args ...any)
}
