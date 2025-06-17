package logger

type Config struct {
	Level          LogLevel
	ComponentLevel map[string]LogLevel
}

func DefaultConfig() Config {
	return Config{
		Level:          LogError,
		ComponentLevel: map[string]LogLevel{},
	}
}
