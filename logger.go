package ultralogger

import (
    "errors"
    "io"
    "os"
)

// Logger defines the interface for a structured ultraLogger in Go.
//
// This interface is useful for either creating your own logger or for using an existing logger, and preventing changes
// to the loggers formatting settings.
type Logger interface {
    // Log logs at the specified level without formatting.
    Log(level Level, data any)

    // Debug logs a debug-level message.
    Debug(data any)

    // Info logs an info-level message.
    Info(data any)

    // Warn logs a warning-level message.
    Warn(data any)

    // Error logs an error-level message.
    Error(data any)

    // Panic logs a panic-level message and then panics.
    Panic(data any)

    // SetMinLevel sets the minimum logging level that will be output.
    SetMinLevel(level Level)

    // SetTag sets the tag for the logger.
    SetTag(tag string)

    Silence(enable bool)
}

var defaultDateTimeFormat = "2006-01-02 15:04:05"
var defaultLevelBracket = BracketAngle

var defaultCurrentTimeField, _ = NewCurrentTimeField("time", defaultDateTimeFormat)

var defaultFields = []Field{
    defaultCurrentTimeField,
    NewLevelField(defaultLevelBracket),
    &FieldMessage{},
}

func NewLoggerWithOptions(opts ...LoggerOption) (Logger, error) {
    l := newUltraLogger()

    for _, opt := range opts {
        if err := opt(l); err != nil {
            return nil, err
        }
    }

    if l.destinations == nil {
        defaultFormatter, _ := NewFormatter(OutputFormatText, defaultFields)
        l.destinations = map[io.Writer]LogLineFormatter{os.Stdout: defaultFormatter}
    }

    return l, nil
}

// NewLogger returns a new Logger that writes to stdout with the default text output format.
func NewLogger() Logger {
    formatter, _ := NewFormatter(OutputFormatText, defaultFields)

    logger, _ := NewLoggerWithOptions(WithStdoutFormatter(formatter))

    return logger
}

//NewFileLogger returns a new Logger that writes to a file.
//
//If the filename is empty, FileNotSpecifiedError is returned.
//If the file does not exist, FileNotFoundError is returned.
func NewFileLogger(filename string, outputFormat OutputFormat) (Logger, error) {
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

    formatter, err := NewFormatter(outputFormat, defaultFields)
    if err != nil {
        return nil, err
    }

    fileLogger, err := NewLoggerWithOptions(WithDestination(filePtr, formatter))
    if err != nil {
        return nil, err
    }

    return fileLogger, nil
}
