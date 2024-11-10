package ultralogger

import (
    "io"
)

type LoggerOption func(l *ultraLogger) error

func WithDestination(writer io.Writer) LoggerOption {
    return func(l *ultraLogger) error {
        l.writer = writer
        return nil
    }
}

func WithMinLevel(level Level) LoggerOption {
    return func(l *ultraLogger) error {
        l.minLevel = level
        return nil
    }
}

func WithFormatter(formatter Formatter) LoggerOption {
    return func(l *ultraLogger) error {
        l.formatter = formatter
        return nil
    }
}

func WithPrefixFields(fields ...Field) LoggerOption {
    return func(l *ultraLogger) error {
        // TODO: Handle this error with a custom type.
        return l.formatter.SetPrefixFields(fields...)
    }
}

func WithSuffixFields(fields ...Field) LoggerOption {
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
        colorizedFormatter, ok := l.formatter.(ColorizedFormatter)
        if !ok {
            return ColorizationNotSupportedError
        }

        return colorizedFormatter.EnableColorization(colorize)
    }
}

func WithLevelColors(colors map[Level]Color) LoggerOption {
    return func(l *ultraLogger) error {
        colorizedFormatter, ok := l.formatter.(ColorizedFormatter)
        if !ok {
            return ColorizationNotSupportedError
        }

        return colorizedFormatter.SetLevelColors(colors)
    }
}
