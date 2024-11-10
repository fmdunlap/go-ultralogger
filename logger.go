package ultralogger

import (
    "errors"
    "github.com/fmdunlap/go-ultralogger/v2/bracket"
    "github.com/fmdunlap/go-ultralogger/v2/field"
    "github.com/fmdunlap/go-ultralogger/v2/formatter"
    "github.com/fmdunlap/go-ultralogger/v2/level"
    "os"
)

// Logger defines the interface for a structured ultraLogger in Go.
//
// This interface is useful for either creating your own logger or for using an existing logger, and preventing changes
// to the loggers formatting settings.
type Logger interface {
    // Log logs at the specified level without formatting.
    Log(level level.Level, msg string)

    // Logf logs at the specified level with formatted message.
    Logf(level level.Level, format string, args ...any)

    // Debug logs a debug-level message.
    Debug(msg string)

    // Debugf logs a debug-level message with formatting.
    Debugf(format string, args ...any)

    // Info logs an info-level message.
    Info(msg string)

    // Infof logs an info-level message with formatting.
    Infof(format string, args ...any)

    // Warn logs a warning-level message.
    Warn(msg string)

    // Warnf logs a warning-level message with formatting.
    Warnf(format string, args ...any)

    // Error logs an error-level message.
    Error(msg string)

    // Errorf logs an error-level message with formatting.
    Errorf(format string, args ...any)

    // Panic logs a panic-level message and then panics.
    Panic(msg string)

    // Panicf logs a panic-level message with formatting and then panics.
    Panicf(format string, args ...any)

    // Slog returns the string representation of a log message with the given level and message.
    Slog(level level.Level, msg string) string

    // Slogf returns the string representation of a formatted log message with the given level and sprint string.
    Slogf(level level.Level, format string, args ...any) string

    // Slogln returns the string representation of a log message with the given level and message, followed by a newline.
    Slogln(level level.Level, msg string) string

    // SetMinLevel sets the minimum logging level that will be output.
    SetMinLevel(level level.Level)
}

var defaultDateTimeFormat = "2006-01-02 15:04:05"
var defaultLevelBracket = bracket.Angle

var defaultPrefixFields = []field.Field{
    field.NewDateTimeField(defaultDateTimeFormat),
    field.NewLevelField(defaultLevelBracket),
}

// NewLogger returns a new Logger that writes to stdout
func NewLogger(opts ...LoggerOption) (Logger, error) {
    fmtr, _ := formatter.NewColorizedFormatter(
        defaultPrefixFields,
        nil,
        false,
    )

    l := &ultraLogger{
        writer:            os.Stdout,
        minLevel:          level.Info,
        formatter:         fmtr,
        silent:            false,
        fallback:          true,
        panicOnPanicLevel: false,
    }

    for _, opt := range opts {
        if err := opt(l); err != nil {
            return nil, err
        }
    }

    return l, nil
}

// NewFileLogger returns a new Logger that writes to a file.
//
// If the filename is empty, FileNotSpecifiedError is returned.
// If the file does not exist, FileNotFoundError is returned.
func NewFileLogger(filename string) (Logger, error) {
    if filename == "" {
        return nil, FileNotSpecifiedError
    }

    var err error
    filePtr, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return nil, &FileNotFoundError{filename: filename}
        }
        return nil, err
    }

    fileLogger, err := NewLogger(WithDestination(filePtr))
    if err != nil {
        return nil, err
    }

    return fileLogger, nil
}
