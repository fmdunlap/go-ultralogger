package ultralogger

import (
    "io"
    "os"
)

type LoggerOption func(l *ultraLogger) error

func WithMinLevel(level Level) LoggerOption {
    return func(l *ultraLogger) error {
        l.minLevel = level
        return nil
    }
}

// WithStdoutFormatter sets the formatter to use for stdout.
// Note: This will not overwrite existing, non-stdout destinations, if any.
func WithStdoutFormatter(formatter LogLineFormatter) LoggerOption {
    return func(l *ultraLogger) error {
        if l.destinations == nil {
            l.destinations = map[io.Writer]LogLineFormatter{}
        }

        l.destinations[os.Stdout] = formatter
        return nil
    }
}

// WithDestination sets the destination for the logger. If the formatter is nil, the destination will be ignored.
// If the logger already has destinations, this will overwrite them.
func WithDestination(destination io.Writer, formatter LogLineFormatter) LoggerOption {
    return func(l *ultraLogger) error {
        l.destinations = map[io.Writer]LogLineFormatter{destination: formatter}
        return nil
    }
}

// WithDestinations sets the destinations for the logger. If the formatter is nil, the destination will be ignored.
// If the logger already has destinations, this will overwrite them.
func WithDestinations(destinations map[io.Writer]LogLineFormatter) LoggerOption {
    return func(l *ultraLogger) error {
        l.destinations = destinations
        return nil
    }
}

func WithSilent(silent bool) LoggerOption {
    return func(l *ultraLogger) error {
        l.silent = silent
        return nil
    }
}

func WithFallbackEnabled(fallback bool) LoggerOption {
    return func(l *ultraLogger) error {
        l.fallback = fallback
        return nil
    }
}

func WithPanicOnPanicLevel(panicOnPanicLevel bool) LoggerOption {
    return func(l *ultraLogger) error {
        l.panicOnPanicLevel = panicOnPanicLevel
        return nil
    }
}

func WithDefaultColorizationEnabled() LoggerOption {
    return func(l *ultraLogger) error {
        if len(l.destinations) == 0 {
            defaultFormatter, _ := NewFormatter(OutputFormatText, defaultFields)
            l.destinations = map[io.Writer]LogLineFormatter{os.Stdout: defaultFormatter}
        }

        l.destinations[os.Stdout] = NewColorizedFormatter(l.destinations[os.Stdout], nil)
        return nil
    }
}

func WithCustomColorization(colors map[Level]Color) LoggerOption {
    return func(l *ultraLogger) error {
        if l.destinations == nil {
            defaultFormatter, _ := NewFormatter(OutputFormatText, defaultFields)
            l.destinations = map[io.Writer]LogLineFormatter{os.Stdout: defaultFormatter}
        }

        l.destinations[os.Stdout] = NewColorizedFormatter(l.destinations[os.Stdout], colors)
        return nil
    }
}

func WithTag(tag string) LoggerOption {
    return func(l *ultraLogger) error {
        l.SetTag(tag)
        return nil
    }
}
