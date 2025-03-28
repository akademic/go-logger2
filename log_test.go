package logger

import (
	"testing"
)

func TestLogString(t *testing.T) {
	testCases := []struct {
		name     string
		log      Log
		expected string
	}{
		{
			name: "Log with component",
			log: Log{
				Level:     LogInfo,
				Component: "TestComponent",
				Message:   "Test message",
			},
			expected: "[inf] [TestComponent]: Test message",
		},
		{
			name: "Log without component",
			log: Log{
				Level:   LogError,
				Message: "Error occurred",
			},
			expected: "[err]: Error occurred",
		},
		{
			name: "Log with empty message",
			log: Log{
				Level:     LogDebug,
				Component: "DebugComponent",
				Message:   "",
			},
			expected: "[dbg] [DebugComponent]: ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.log.String()
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestLogLabels(t *testing.T) {
	testCases := []struct {
		name     string
		log      Log
		expected map[string]string
	}{
		{
			name: "Log with component",
			log: Log{
				Level:     LogInfo,
				Component: "TestComponent",
				Message:   "Test message",
			},
			expected: map[string]string{
				"component": "TestComponent",
			},
		},
		{
			name: "Log without component",
			log: Log{
				Level:   LogError,
				Message: "Error occurred",
			},
			expected: nil,
		},
		{
			name: "Log with empty component",
			log: Log{
				Level:     LogDebug,
				Component: "",
				Message:   "Debug message",
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.log.Labels()

			if len(result) != len(tc.expected) {
				t.Errorf("Expected labels length %d, got %d", len(tc.expected), len(result))
				return
			}

			if tc.expected == nil {
				if result != nil {
					t.Errorf("Expected nil labels, got %v", result)
				}
				return
			}

			for k, v := range tc.expected {
				if resultV, exists := result[k]; !exists || resultV != v {
					t.Errorf("Expected label %s=%s, got %v", k, v, result)
				}
			}
		})
	}
}
