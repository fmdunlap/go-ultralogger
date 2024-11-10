package ultralogger

import (
    "fmt"
    "github.com/fmdunlap/go-ultralogger/formatter"
    "github.com/fmdunlap/go-ultralogger/level"
    "io"
    "os"
)

type ultraLogger struct {
    writer            io.Writer
    minLevel          level.Level
    formatter         formatter.Formatter
    silent            bool
    fallback          bool
    panicOnPanicLevel bool
}

// Log logs a message with the given level and message.
func (l *ultraLogger) Log(level level.Level, msg string) {
    outBytes := []byte(l.Slogln(level, msg))

    if _, err := l.writer.Write(outBytes); err != nil {
        l.handleLogWriterError(level, msg, err)
    }
}

// Logf logs a formatted message with the given level and sprint string.
func (l *ultraLogger) Logf(level level.Level, format string, args ...any) {
    l.Log(level, fmt.Sprintf(format, args...))
}

// Debug logs a message with the Debug level and message.
func (l *ultraLogger) Debug(msg string) {
    l.Log(level.Debug, msg)
}

// Debugf logs a formatted message with the Debug level and sprint string.
func (l *ultraLogger) Debugf(format string, args ...any) {
    l.Logf(level.Debug, format, args...)
}

// Info logs a message with the Info level and message.
func (l *ultraLogger) Info(msg string) {
    l.Log(level.Info, msg)
}

// Infof logs a formatted message with the Info level and sprint string.
func (l *ultraLogger) Infof(format string, args ...any) {
    l.Logf(level.Info, format, args...)
}

// Warn logs a message with the Warn level and message.
func (l *ultraLogger) Warn(msg string) {
    l.Log(level.Warn, msg)
}

// Warnf logs a formatted message with the Warn level and sprint string.
func (l *ultraLogger) Warnf(format string, args ...any) {
    l.Logf(level.Warn, format, args...)
}

// Error logs a message with the Error level and message.
func (l *ultraLogger) Error(msg string) {
    l.Log(level.Error, msg)
}

// Errorf logs a formatted message with the Error level and sprint string.
func (l *ultraLogger) Errorf(format string, args ...any) {
    l.Logf(level.Error, format, args...)
}

// Panic logs a message with the Panic level and message. If panicOnPanicLevel is true, it panics.
func (l *ultraLogger) Panic(msg string) {
    l.Log(level.Panic, msg)

    if l.panicOnPanicLevel {
        panic(msg)
    }
}

// Panicf logs a formatted message with the Panic level and sprint string. If panicOnPanicLevel is true, it panics.
func (l *ultraLogger) Panicf(format string, args ...any) {
    l.Logf(level.Panic, format, args...)

    if l.panicOnPanicLevel {
        panic(fmt.Sprintf(format, args...))
    }
}

// Slog returns the string representation of a log message with the given level and message.
func (l *ultraLogger) Slog(level level.Level, msg string) string {
    if l.silent || level < l.minLevel {
        return ""
    }

    return l.formatter.Format(level, msg)
}

// Slogf returns the string representation of a formatted log message with the given level and sprint string.
func (l *ultraLogger) Slogf(level level.Level, format string, args ...any) string {
    // Optimize-out the Sprintf if the level is too low or silent.
    if l.silent || level < l.minLevel {
        return ""
    }

    return l.formatter.Formatf(level, format, args...)
}

// Slogln returns the string representation of a log message with the given level and message, followed by a newline.
func (l *ultraLogger) Slogln(level level.Level, msg string) string {
    // Optimize-out the Sprintf if the level is too low or silent.
    if l.silent || level < l.minLevel {
        return ""
    }

    return l.formatter.Format(level, msg) + "\n"
}

func (l *ultraLogger) SetMinLevel(level level.Level) {
    l.minLevel = level
}

// handleLogWriterError handles errors that occur while writing to the output. On failure, the log will fall back to
// writing to os.Stdout.
func (l *ultraLogger) handleLogWriterError(msgLevel level.Level, msg string, err error) {
    if !l.fallback || l.writer == os.Stdout {
        panic(err)
    }

    l.writer = os.Stdout
    l.Logf(
        level.Error,
        "error writing to original log writer, falling back to stdout: %v",
        err,
    )
    l.Log(msgLevel, msg)
}
