package ultralogger

import (
    "context"
    "fmt"
    "io"
    "os"
    "time"
)

const loglineTimeout = time.Millisecond * 250

type ultraLogger struct {
    minLevel          Level
    destinations      map[io.Writer]LogLineFormatter
    tag               string
    silent            bool
    fallback          bool
    panicOnPanicLevel bool
    async             bool
}

func newUltraLogger() *ultraLogger {
    return &ultraLogger{
        minLevel:          Info,
        destinations:      map[io.Writer]LogLineFormatter{},
        silent:            false,
        fallback:          true,
        panicOnPanicLevel: false,
        async:             true,
    }
}

// Log logs a message with the given level and message.
func (l *ultraLogger) Log(level Level, data any) {
    if l.silent || level < l.minLevel {
        return
    }

    args := LogLineArgs{
        Level: level,
        Tag:   l.tag,
    }

    for w, f := range l.destinations {
        if f == nil {
            continue
        }

        if l.async {
            go l.writeLogLineAsync(w, f, args, data, loglineTimeout)
            continue
        }

        l.writeLogLine(w, f, args, data)
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

func (l *ultraLogger) writeLogLine(
    w io.Writer,
    f LogLineFormatter,
    args LogLineArgs,
    data any,
) {
    formatResult := f.FormatLogLine(args, data)
    if formatResult.err != nil {
        l.Error(fmt.Sprintf("failed to format log line. formatter=%v, data=%v, err=%v", f, data, formatResult.err))
        return
    }

    writeResult := write(w, formatResult.bytes)
    if writeResult != nil {
        l.handleLogWriterError(w, args.Level, data, writeResult)
    }
}

func (l *ultraLogger) writeLogLineAsync(
    w io.Writer,
    f LogLineFormatter,
    args LogLineArgs,
    data any,
    timeout time.Duration,
) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    fmtChan := make(chan FormatResult, 1)
    go formatLogLineAsync(ctx, fmtChan, args, f, data)

    var logBytes []byte
    select {
    case result := <-fmtChan:
        if result.err != nil {
            l.Error(fmt.Sprintf("failed to format log line. formatter=%v, data=%v, err=%v", f, data, result.err))
            return
        }

        if len(result.bytes) == 0 {
            return
        }

        logBytes = result.bytes
    case <-ctx.Done():
        return
    }

    writeChan := make(chan error, 1)
    go writeLogLineAsync(ctx, writeChan, w, logBytes)

    select {
    case err := <-writeChan:
        if err != nil {
            l.handleLogWriterError(w, args.Level, data, err)
        }
    case <-ctx.Done():
        return
    }
}

func formatLogLineAsync(
    ctx context.Context,
    resultChan chan FormatResult,
    args LogLineArgs,
    formatter LogLineFormatter,
    data any,
) {
    defer close(resultChan)

    select {
    case <-ctx.Done():
        return
    case resultChan <- formatter.FormatLogLine(args, data):
    }
}

func writeLogLineAsync(
    ctx context.Context,
    resultChan chan error,
    w io.Writer,
    b []byte,
) {
    defer close(resultChan)

    select {
    case <-ctx.Done():
        return
    case resultChan <- write(w, b):
    }
}

func write(w io.Writer, b []byte) error {
    _, err := w.Write(append(b, '\n'))
    return err
}
