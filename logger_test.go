package logger

import (
	"maps"
	"testing"
)

type mockBaseLogger struct {
	printCalls [][]any
}

func (m *mockBaseLogger) Print(v ...any) {
	m.printCalls = append(m.printCalls, v)
}

func TestNew(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		component string
	}{
		{
			name:      "with nil config should use default",
			config:    nil,
			component: "test-component",
		},
		{
			name: "with custom config",
			config: &Config{
				Level:          LogInfo,
				ComponentLevel: map[string]LogLevel{"test-component": LogDebug},
			},
			component: "test-component",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseLogger := &mockBaseLogger{}

			logger := New(baseLogger, tt.component, tt.config)

			if logger == nil {
				t.Fatal("expected logger to not be nil")
			}

			if logger.component != tt.component {
				t.Errorf("expected component %q, got %q", tt.component, logger.component)
			}

			if logger.baseLogger != baseLogger {
				t.Error("expected baseLogger to be set correctly")
			}

			if tt.config == nil {
				// Should use default config
				if logger.config.Level != LogError {
					t.Errorf("expected default level LogError, got %v", logger.config.Level)
				}
				if logger.config.ComponentLevel == nil {
					t.Error("expected ComponentLevel to be initialized")
				}
			} else {
				if logger.config.Level != tt.config.Level {
					t.Errorf("expected level %v, got %v", tt.config.Level, logger.config.Level)
				}
			}
		})
	}
}

func TestLoggerImpl_WithComponent(t *testing.T) {
	baseLogger := &mockBaseLogger{}
	config := &Config{Level: LogInfo}
	originalLogger := New(baseLogger, "original", config)

	newLogger := originalLogger.WithComponent("new-component")

	if newLogger == originalLogger {
		t.Error("expected new logger instance")
	}

	newLoggerImpl, ok := newLogger.(*LoggerImpl)
	if !ok {
		t.Fatal("expected *LoggerImpl")
	}

	if newLoggerImpl.component != "new-component" {
		t.Errorf("expected component 'new-component', got %q", newLoggerImpl.component)
	}

	if newLoggerImpl.config != originalLogger.config {
		t.Error("expected same config reference")
	}

	if newLoggerImpl.baseLogger != originalLogger.baseLogger {
		t.Error("expected same baseLogger reference")
	}
}

func TestLoggerImpl_Info(t *testing.T) {
	tests := []struct {
		name          string
		config        *Config
		component     string
		pattern       string
		args          []any
		expectLog     bool
		expectedMsg   string
		expectedLevel LogLevel
	}{
		{
			name:          "info enabled globally",
			config:        &Config{Level: LogInfo, ComponentLevel: map[string]LogLevel{}},
			component:     "",
			pattern:       "test message",
			args:          nil,
			expectLog:     true,
			expectedMsg:   "test message",
			expectedLevel: LogInfo,
		},
		{
			name:      "info disabled globally",
			config:    &Config{Level: LogError, ComponentLevel: map[string]LogLevel{}},
			component: "",
			pattern:   "test message",
			args:      nil,
			expectLog: false,
		},
		{
			name:          "info with formatting",
			config:        &Config{Level: LogInfo, ComponentLevel: map[string]LogLevel{}},
			component:     "test",
			pattern:       "user %s has %d items",
			args:          []any{"john", 5},
			expectLog:     true,
			expectedMsg:   "user john has 5 items",
			expectedLevel: LogInfo,
		},
		{
			name:          "component specific level allows info",
			config:        &Config{Level: LogError, ComponentLevel: map[string]LogLevel{"test": LogInfo}},
			component:     "test",
			pattern:       "test message",
			args:          nil,
			expectLog:     true,
			expectedMsg:   "test message",
			expectedLevel: LogInfo,
		},
		{
			name:      "component specific level blocks info",
			config:    &Config{Level: LogDebug, ComponentLevel: map[string]LogLevel{"test": LogError}},
			component: "test",
			pattern:   "test message",
			args:      nil,
			expectLog: false,
		},
		{
			name:          "no info with nil config",
			config:        nil,
			component:     "test",
			pattern:       "test message",
			args:          nil,
			expectLog:     false,
			expectedMsg:   "",
			expectedLevel: LogInfo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseLogger := &mockBaseLogger{}
			logger := New(baseLogger, tt.component, tt.config)

			logger.Info(tt.pattern, tt.args...)

			if tt.expectLog {
				if len(baseLogger.printCalls) != 1 {
					t.Fatalf("expected 1 print call, got %d", len(baseLogger.printCalls))
				}

				if len(baseLogger.printCalls[0]) != 1 {
					t.Fatalf("expected 1 argument to Print, got %d", len(baseLogger.printCalls[0]))
				}

				log, ok := baseLogger.printCalls[0][0].(Log)
				if !ok {
					t.Fatalf("expected Log type, got %T", baseLogger.printCalls[0][0])
				}

				if log.Level != tt.expectedLevel {
					t.Errorf("expected level %v, got %v", tt.expectedLevel, log.Level)
				}

				if log.Message != tt.expectedMsg {
					t.Errorf("expected message %q, got %q", tt.expectedMsg, log.Message)
				}

				if log.Component != tt.component {
					t.Errorf("expected component %q, got %q", tt.component, log.Component)
				}
			} else {
				if len(baseLogger.printCalls) != 0 {
					t.Errorf("expected no print calls, got %d", len(baseLogger.printCalls))
				}
			}
		})
	}
}

func TestLoggerImpl_Debug(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		component   string
		pattern     string
		args        []any
		expectLog   bool
		expectedMsg string
	}{
		{
			name:        "debug enabled",
			config:      &Config{Level: LogDebug, ComponentLevel: map[string]LogLevel{}},
			component:   "",
			pattern:     "debug message",
			args:        nil,
			expectLog:   true,
			expectedMsg: "debug message",
		},
		{
			name:      "debug disabled",
			config:    &Config{Level: LogError, ComponentLevel: map[string]LogLevel{}},
			component: "",
			pattern:   "debug message",
			args:      nil,
			expectLog: false,
		},
		{
			name:        "debug with formatting",
			config:      &Config{Level: LogDebug, ComponentLevel: map[string]LogLevel{}},
			component:   "debugger",
			pattern:     "value is %v",
			args:        []any{42},
			expectLog:   true,
			expectedMsg: "value is 42",
		},
		{
			name:        "debug with nil config",
			config:      nil,
			component:   "debugger",
			pattern:     "debug message",
			args:        nil,
			expectLog:   false,
			expectedMsg: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseLogger := &mockBaseLogger{}
			logger := New(baseLogger, tt.component, tt.config)

			logger.Debug(tt.pattern, tt.args...)

			if tt.expectLog {
				if len(baseLogger.printCalls) != 1 {
					t.Fatalf("expected 1 print call, got %d", len(baseLogger.printCalls))
				}

				log, ok := baseLogger.printCalls[0][0].(Log)
				if !ok {
					t.Fatalf("expected Log type, got %T", baseLogger.printCalls[0][0])
				}

				if log.Level != LogDebug {
					t.Errorf("expected level LogDebug, got %v", log.Level)
				}

				if log.Message != tt.expectedMsg {
					t.Errorf("expected message %q, got %q", tt.expectedMsg, log.Message)
				}
			} else {
				if len(baseLogger.printCalls) != 0 {
					t.Errorf("expected no print calls, got %d", len(baseLogger.printCalls))
				}
			}
		})
	}
}

func TestLoggerImpl_Error(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		component   string
		pattern     string
		args        []any
		expectLog   bool
		expectedMsg string
	}{
		{
			name:        "error always logged when level allows",
			config:      &Config{Level: LogError, ComponentLevel: map[string]LogLevel{}},
			component:   "",
			pattern:     "error occurred",
			args:        nil,
			expectLog:   true,
			expectedMsg: "error occurred",
		},
		{
			name:      "error disabled when level is off",
			config:    &Config{Level: LogOff, ComponentLevel: map[string]LogLevel{}},
			component: "",
			pattern:   "error occurred",
			args:      nil,
			expectLog: false,
		},
		{
			name:        "error with formatting",
			config:      &Config{Level: LogError, ComponentLevel: map[string]LogLevel{}},
			component:   "error-handler",
			pattern:     "failed to process %s: %v",
			args:        []any{"file.txt", "permission denied"},
			expectLog:   true,
			expectedMsg: "failed to process file.txt: permission denied",
		},
		{
			name:        "error with nil config",
			config:      nil,
			component:   "error-handler",
			pattern:     "error occurred",
			args:        nil,
			expectLog:   true,
			expectedMsg: "error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseLogger := &mockBaseLogger{}
			logger := New(baseLogger, tt.component, tt.config)

			logger.Error(tt.pattern, tt.args...)

			if tt.expectLog {
				if len(baseLogger.printCalls) != 1 {
					t.Fatalf("expected 1 print call, got %d", len(baseLogger.printCalls))
				}

				log, ok := baseLogger.printCalls[0][0].(Log)
				if !ok {
					t.Fatalf("expected Log type, got %T", baseLogger.printCalls[0][0])
				}

				if log.Level != LogError {
					t.Errorf("expected level LogError, got %v", log.Level)
				}

				if log.Message != tt.expectedMsg {
					t.Errorf("expected message %q, got %q", tt.expectedMsg, log.Message)
				}
			} else {
				if len(baseLogger.printCalls) != 0 {
					t.Errorf("expected no print calls, got %d", len(baseLogger.printCalls))
				}
			}
		})
	}
}

func TestLoggerImpl_SetConfig(t *testing.T) {
	tests := []struct {
		name          string
		initialConfig *Config
		newConfig     *Config
		expectChange  bool
	}{
		{
			name:          "set nil config should not change",
			initialConfig: &Config{Level: LogError, ComponentLevel: map[string]LogLevel{}},
			newConfig:     nil,
			expectChange:  false,
		},
		{
			name:          "set new level",
			initialConfig: &Config{Level: LogError, ComponentLevel: map[string]LogLevel{}},
			newConfig:     &Config{Level: LogDebug, ComponentLevel: nil},
			expectChange:  true,
		},
		{
			name:          "set new component levels",
			initialConfig: &Config{Level: LogError, ComponentLevel: map[string]LogLevel{}},
			newConfig:     &Config{Level: LogError, ComponentLevel: map[string]LogLevel{"test": LogInfo}},
			expectChange:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseLogger := &mockBaseLogger{}
			logger := New(baseLogger, "test", tt.initialConfig)

			originalLevel := logger.config.Level
			originalComponentLevel := logger.config.ComponentLevel

			logger.SetConfig(tt.newConfig)

			if !tt.expectChange {
				if logger.config.Level != originalLevel {
					t.Errorf("expected level to remain %v, got %v", originalLevel, logger.config.Level)
				}
				return
			}

			if tt.newConfig.Level != "" && logger.config.Level != tt.newConfig.Level {
				t.Errorf("expected level %v, got %v", tt.newConfig.Level, logger.config.Level)
			}

			if tt.newConfig.ComponentLevel != nil {
				if maps.Equal(logger.config.ComponentLevel, originalComponentLevel) {
					t.Error("expected ComponentLevel to be updated")
				}
			}
		})
	}
}

func TestLoggerImpl_logOn(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		component string
		level     LogLevel
		expected  bool
	}{
		{
			name:      "no component uses global level",
			config:    &Config{Level: LogInfo, ComponentLevel: map[string]LogLevel{}},
			component: "",
			level:     LogInfo,
			expected:  true,
		},
		{
			name:      "component not in map uses global level",
			config:    &Config{Level: LogError, ComponentLevel: map[string]LogLevel{"other": LogDebug}},
			component: "missing",
			level:     LogError,
			expected:  true,
		},
		{
			name:      "component uses specific level",
			config:    &Config{Level: LogError, ComponentLevel: map[string]LogLevel{"test": LogDebug}},
			component: "test",
			level:     LogInfo,
			expected:  true,
		},
		{
			name:      "component specific level blocks log",
			config:    &Config{Level: LogDebug, ComponentLevel: map[string]LogLevel{"test": LogError}},
			component: "test",
			level:     LogInfo,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseLogger := &mockBaseLogger{}
			logger := New(baseLogger, tt.component, tt.config)

			result := logger.logOn(tt.level)

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
