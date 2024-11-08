package ultralogger

import (
    "errors"
    "os"
)

// Logger defines the interface for a structured UltraLogger in Go.
//
// This interface is useful for either creating your own logger or for using an existing logger, and preventing changes
// to the loggers formatting settings.
type Logger interface {
    // Log logs at the specified level without formatting.
    Log(level Level, msg string)

    // Logf logs at the specified level with formatted message.
    Logf(level Level, format string, args ...any)

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

    // SetMinLevel sets the minimum logging level that will be output.
    SetMinLevel(level Level) Logger
}

// NewStdoutLogger returns a new Logger that writes to stdout
func NewStdoutLogger() *UltraLogger {
    l := NewUltraLogger(os.Stdout)
    return l
}

// NewFileLogger returns a new Logger that writes to a file.
//
// If the filename is empty, FileNotSpecifiedError is returned.
// If the file does not exist, FileNotFoundError is returned.
func NewFileLogger(filename string) (*UltraLogger, error) {

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

    fileLogger := NewUltraLogger(filePtr)
    fileLogger.colorize = false
    return fileLogger, nil
}
