package ultralogger

import (
    "fmt"
    "io"
    "os"
)

type ultraLogger struct {
    minLevel          Level
    destinations      map[io.Writer]LogLineFormatter
    tag               string
    silent            bool
    fallback          bool
    panicOnPanicLevel bool
}

func newUltraLogger() *ultraLogger {
    return &ultraLogger{
        minLevel:          Info,
        destinations:      map[io.Writer]LogLineFormatter{},
        silent:            false,
        fallback:          true,
        panicOnPanicLevel: false,
    }
}

func (l *ultraLogger) writeLogLine(writer io.Writer, formatter LogLineFormatter, logLineContext LogLineContext, data any) {
    if formatter == nil {
        return
    }

    outBytes, err := formatter.FormatLogLine(logLineContext, data)

    if err != nil {
        l.Error(fmt.Sprintf("failed to format log line. formatter=%v, data=%v, err=%v", formatter, data, err))
        return
    }

    if len(outBytes) == 0 {
        return
    }

    if _, err := writer.Write(fmt.Append(outBytes, "\n")); err != nil {
        l.handleLogWriterError(writer, logLineContext.Level, data, err)
    }
}

// Log logs a message with the given level and message.
func (l *ultraLogger) Log(level Level, data any) {
    if l.silent || level < l.minLevel {
        return
    }

    loglineContext := LogLineContext{
        Level: level,
        Tag:   l.tag,
    }

    for writer, formatter := range l.destinations {
        go l.writeLogLine(writer, formatter, loglineContext, data)
    }
}

// Debug logs a message with the Debug level and message.
func (l *ultraLogger) Debug(data any) {
    l.Log(Debug, data)
}

// Info logs a message with the Info level and message.
func (l *ultraLogger) Info(data any) {
    l.Log(Info, data)
}

// Warn logs a message with the Warn level and message.
func (l *ultraLogger) Warn(data any) {
    l.Log(Warn, data)
}

// Error logs a message with the Error level and message.
func (l *ultraLogger) Error(data any) {
    l.Log(Error, data)
}

// Panic logs a message with the Panic level and message. If panicOnPanicLevel is true, it panics.
func (l *ultraLogger) Panic(data any) {
    l.Log(Panic, data)

    if l.panicOnPanicLevel {
        panic(data)
    }
}

func (l *ultraLogger) SetMinLevel(level Level) {
    l.minLevel = level
}

func (l *ultraLogger) SetTag(tag string) {
    l.tag = tag
}

func (l *ultraLogger) Silence(enable bool) {
    l.silent = enable
}

// handleLogWriterError handles errors that occur while writing to the output. On failure, the log will fall back to
// writing to os.Stdout.
func (l *ultraLogger) handleLogWriterError(writer io.Writer, msgLevel Level, msg any, err error) {
    if !l.fallback || writer == os.Stdout {
        panic(err)
    }

    l.destinations[writer] = nil
    l.Error(
        fmt.Sprintf("error writing to original log writer, disabling formatter for writer: %v", err),
    )
    l.Log(msgLevel, msg)
}
