package ultralogger

import (
    "github.com/fmdunlap/go-ultralogger/color"
    "github.com/fmdunlap/go-ultralogger/field"
    "github.com/fmdunlap/go-ultralogger/formatter"
    "github.com/fmdunlap/go-ultralogger/level"
    "io"
)

type LoggerOption func(l *ultraLogger) error

func WithDestination(writer io.Writer) LoggerOption {
    return func(l *ultraLogger) error {
        l.writer = writer
        return nil
    }
}

func WithMinLevel(level level.Level) LoggerOption {
    return func(l *ultraLogger) error {
        l.minLevel = level
        return nil
    }
}

func WithFormatter(formatter formatter.Formatter) LoggerOption {
    return func(l *ultraLogger) error {
        l.formatter = formatter
        return nil
    }
}

func WithPrefixFields(fields ...field.Field) LoggerOption {
    return func(l *ultraLogger) error {
        // TODO: Handle this error with a custom type.
        return l.formatter.SetPrefixFields(fields...)
    }
}

func WithSuffixFields(fields ...field.Field) LoggerOption {
    return func(l *ultraLogger) error {
        // TODO: Handle this error with a custom type.
        return l.formatter.SetSuffixFields(fields...)
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

func WithColorization(colorize bool) LoggerOption {
    return func(l *ultraLogger) error {
        colorizedFormatter, ok := l.formatter.(formatter.ColorizedFormatter)
        if !ok {
            return ColorizationNotSupportedError
        }

        return colorizedFormatter.EnableColorization(colorize)
    }
}

func WithLevelColors(colors map[level.Level]color.Color) LoggerOption {
    return func(l *ultraLogger) error {
        colorizedFormatter, ok := l.formatter.(formatter.ColorizedFormatter)
        if !ok {
            return ColorizationNotSupportedError
        }

        return colorizedFormatter.SetLevelColors(colors)
    }
}
