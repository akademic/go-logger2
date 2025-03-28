## What is it for?

This package is a simple logging library for Go.
It provides a flexible logging interface with support for different log levels and component-specific logging configurations.

## FEATURES

- Only stdlib dependencies
- Error, Info, Debug levels (Errors print always, Info and Debug only if configured)
- Logging with prefixes. Log output with particular prefix can be controlled with config with same levels (Error, Info, Debug)

## INSTALL

```
go get github.com/akademic/go-logger2
```

## USAGE

### Basic Example

```go
package main

import (
    "log"
    logger "github.com/akademic/go-logger2"
)

func main() {
    // Create a base logger
    baseLogger := log.New(os.Stdout, "", log.LstdFlags)

    // Create a config
    config := &logger.Config{
        Level: logger.LogDebug,
        ComponentLevel: map[string]logger.LogLevel{
            "db":  logger.LogError,
            "api": logger.LogInfo,
        },
    }

    // Create a new logger
    l := logger.New(baseLogger, "", config)

    l.Error("number: %d", 1)   // Will print error log
    l.Info("number: %d", 2)    // Will print info log
    l.Debug("number: %d", 3)   // Will print debug log
}
```

### Component-Specific Logging

```go
package main

import (
    "log"
    logger "github.com/akademic/go-logger2"
)

func main() {
    baseLogger := log.New(os.Stdout, "", log.LstdFlags)
    config := &logger.Config{
        Level: logger.LogDebug,
        ComponentLevel: map[string]logger.LogLevel{
            "db":  logger.LogError,
            "api": logger.LogInfo,
        },
    }

    // Create a logger
    l := logger.New(baseLogger, "", config)

    // Create component-specific loggers
    dbLogger := l.WithComponent("db")
    apiLogger := l.WithComponent("api")

    dbLogger.Error("database error: %v", "connection failed")   // Will print
    dbLogger.Info("database info")                              // Will NOT print (db level is Error)
    dbLogger.Debug("database debug")                            // Will NOT print (db level is Error)

    apiLogger.Error("api error: %v", "validation failed")       // Will print
    apiLogger.Info("api started")                               // Will print
    apiLogger.Debug("api debug details")                        // Will NOT print (api level is Info)
}
```

### Dynamically Changing Configuration

```go
package main

import (
    "log"
    logger "github.com/akademic/go-logger2"
)

func main() {
    baseLogger := log.New(os.Stdout, "", log.LstdFlags)
    config := &logger.Config{
        Level: logger.LogDebug,
    }

    l := logger.New(baseLogger, "", config)

    l.Info("initial log")  // Will print

    // Update configuration
    l.SetConfig(&logger.Config{
        Level: logger.LogError,
    })

    l.Info("this will NOT print")  // Will NOT print due to new config
}
```

### Writing log to file

```go
package main

import (
    "log"
    "os"
    logger "github.com/akademic/go-logger2"
)

func main() {
    // Create a file
    f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()

    // Create a base logger
    baseLogger := log.New(f, "", log.LstdFlags)

    // Create a logger
    l := logger.New(baseLogger, "", &logger.Config{
        Level: logger.LogDebug,
    })

    l.Info("log to file")
}
```

### Writing log to loki

```go
package main

import (
    "log"
    "os"
    logger "github.com/akademic/go-logger2"
    loki "github.com/akademic/go-logger2-loki"
)

func main() {
    lokiLogger := loki.New(loki.Config{
        Address: "http://localhost:3100",
        Timeout: 500 * time.Millisecond,
        Labels: map[string]string{
            "project": "my-project",
        },
    })

    // Create a logger
    l := logger.New(lokiLogger, "api", &logger.Config{
        Level: logger.LogDebug,
    })

    l.Info("log to loki") // Will send log to loki with labels "project=my-project" and "component=api"
}
```

### Writing log to multiple outputs for io.Writer

```go
package main

import (
    "log"
    "os"
    logger "github.com/akademic/go-logger2"
)

func main() {
    // Create a file
    f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()

    // Create a base logger
    baseLogger := log.New(io.MultiWriter(os.Stdout, f), "", log.LstdFlags)

    // Create a logger
    l := logger.New(baseLogger, "", &logger.Config{
        Level: logger.LogDebug,
    })

    l.Info("log to stdout and file")
}
```

### Writing log to multiple outputs for BaseLogger

```go
package main

import (
    "log"
    "os"
    logger "github.com/akademic/go-logger2"
    loki "github.com/akademic/go-logger2-loki"
)

func main() {
    // Create a file
    f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()

    fileLogger := log.New(f, "", log.LstdFlags)

    stdoutLogger := log.New(os.Stdout, "", log.LstdFlags)

    lokiLogger := loki.New(loki.Config{
        Address: "http://localhost:3100",
        Timeout: 500 * time.Millisecond,
        Labels: map[string]string{
            "project": "my-project",
        },
    })

    multiLogger := logger.NewMultiLogger(fileLogger, stdoutLogger)

    // Create a logger
    l := logger.New(multiLogger, "", &logger.Config{
        Level: logger.LogDebug,
    })

    l.Info("log to stdout and file and loki")
}
```

### Interface

The logger provides a simple interface:

```go
type Logger interface {
    WithComponent(component string) Logger
    Error(format string, args ...any)
    Info(format string, args ...any)
    Debug(format string, args ...any)
}
```

## Important Notes

- Implement your own `BaseLogger` interface to customize log output
- You can use any io.Writer as log output with log.New()
- Log levels can be configured globally and per component
- Error logs are always logged unless explicitly disabled
- Info and Debug logs are configurable

