package ultralogger

import (
    "testing"
)

// TODO: Tests for WithFallback, WithDestination, WithFormatter, WithPanicOnPanicLevel

func testFormatOption(t *testing.T, options []LoggerOption, msg string, msgLevel Level, want string) {
    allOpt := []LoggerOption{WithPrefixFields()}
    allOpt = append(allOpt, options...)

    logger, err := NewLogger(allOpt...)
    if err != nil {
        t.Errorf("NewLogger() error = %v, unexpected", err)
    }

    if got := logger.Slog(msgLevel, msg); got != want {
        t.Errorf("NewLogger() got = %v, want %v", got, want)
    }
}

func TestWithColorization(t *testing.T) {
    tests := []struct {
        name     string
        options  []LoggerOption
        msg      string
        msgLevel Level
        want     string
    }{
        {
            name:     "WithColorization(true)",
            options:  []LoggerOption{WithColorization(true)},
            msg:      "test",
            msgLevel: Warn,
            want:     ColorYellow.Colorize("test"),
        },
        {
            name:     "WithColorization(false)",
            options:  []LoggerOption{WithColorization(false)},
            msg:      "test",
            msgLevel: Warn,
            want:     "test",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            testFormatOption(t, tt.options, tt.msg, tt.msgLevel, tt.want)
        })
    }
}

func TestWithLevelColors(t *testing.T) {
    tests := []struct {
        name     string
        options  []LoggerOption
        msg      string
        msgLevel Level
        want     string
    }{
        {
            name: "WithLevelColors()",
            options: []LoggerOption{WithLevelColors(map[Level]Color{
                Debug: ColorYellow,
                Info:  ColorBlue,
                Warn:  ColorGreen,
                Error: ColorRed,
                Panic: ColorMagenta,
            }),
                WithColorization(true)},
            msg:      "test",
            msgLevel: Warn,
            want:     ColorGreen.Colorize("test"),
        },
        {
            name: "With partial level colors",
            options: []LoggerOption{WithLevelColors(map[Level]Color{
                Debug: ColorYellow,
                Info:  ColorBlue,
                Warn:  ColorGreen,
            }),
                WithColorization(true)},
            msg:      "test",
            msgLevel: Warn,
            want:     ColorGreen.Colorize("test"),
        },
        {
            name: "With partial level colors gets default",
            options: []LoggerOption{WithLevelColors(map[Level]Color{
                Debug: ColorYellow,
                Info:  ColorBlue,
                Warn:  ColorGreen,
            }),
                WithColorization(true)},
            msg:      "test",
            msgLevel: Error,
            want:     ColorRed.Colorize("test"),
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            testFormatOption(t, tt.options, tt.msg, tt.msgLevel, tt.want)
        })
    }
}

func TestWithMinLevel(t *testing.T) {
    tests := []struct {
        name     string
        options  []LoggerOption
        msg      string
        msgLevel Level
        want     string
    }{
        {
            name:     "WithMinLevel(level.Debug), debug shows up",
            options:  []LoggerOption{WithMinLevel(Debug)},
            msg:      "test",
            msgLevel: Debug,
            want:     "test",
        },
        {
            name:     "WithMinLevel(level.Info), debug doesn't show up",
            options:  []LoggerOption{WithMinLevel(Info)},
            msg:      "test",
            msgLevel: Debug,
            want:     "",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            testFormatOption(t, tt.options, tt.msg, tt.msgLevel, tt.want)
        })
    }
}

func TestWithPrefixFields(t *testing.T) {
    tests := []struct {
        name     string
        options  []LoggerOption
        msg      string
        msgLevel Level
        want     string
    }{
        {
            name:     "WithPrefixFields()",
            options:  []LoggerOption{WithPrefixFields(NewLevelField(BracketAngle))},
            msg:      "test",
            msgLevel: Warn,
            want:     "<WARN> test",
        },
        {
            name:     "Multiple Prefix Fields",
            options:  []LoggerOption{WithPrefixFields(NewLevelField(BracketAngle), NewTagField("tag"))},
            msg:      "test",
            msgLevel: Warn,
            want:     "<WARN> [tag] test",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            testFormatOption(t, tt.options, tt.msg, tt.msgLevel, tt.want)
        })
    }
}

func TestWithSilent(t *testing.T) {
    tests := []struct {
        name     string
        options  []LoggerOption
        msg      string
        msgLevel Level
        want     string
    }{
        {
            name:     "WithSilent(true), error doesn't show up",
            options:  []LoggerOption{WithSilent(true)},
            msg:      "test",
            msgLevel: Error,
            want:     "",
        },
        {
            name:     "WithSilent(false), error shows up",
            options:  []LoggerOption{WithSilent(false)},
            msg:      "test",
            msgLevel: Error,
            want:     "test",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            testFormatOption(t, tt.options, tt.msg, tt.msgLevel, tt.want)
        })
    }
}

func TestWithSuffixFields(t *testing.T) {
    tests := []struct {
        name     string
        options  []LoggerOption
        msg      string
        msgLevel Level
        want     string
    }{
        {
            name:     "WithPrefixFields()",
            options:  []LoggerOption{WithSuffixFields(NewLevelField(BracketAngle))},
            msg:      "test",
            msgLevel: Warn,
            want:     "test <WARN>",
        },
        {
            name:     "Multiple Prefix Fields",
            options:  []LoggerOption{WithSuffixFields(NewLevelField(BracketAngle), NewTagField("tag"))},
            msg:      "test",
            msgLevel: Warn,
            want:     "test <WARN> [tag]",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            testFormatOption(t, tt.options, tt.msg, tt.msgLevel, tt.want)
        })
    }
}
