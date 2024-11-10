package ultralogger

import (
    "github.com/fmdunlap/go-ultralogger/v2/bracket"
    "github.com/fmdunlap/go-ultralogger/v2/color"
    "github.com/fmdunlap/go-ultralogger/v2/field"
    "github.com/fmdunlap/go-ultralogger/v2/level"
    "testing"
)

// TODO: Tests for WithFallback, WithDestination, WithFormatter, WithPanicOnPanicLevel

func testFormatOption(t *testing.T, options []LoggerOption, msg string, msgLevel level.Level, want string) {
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
        msgLevel level.Level
        want     string
    }{
        {
            name:     "WithColorization(true)",
            options:  []LoggerOption{WithColorization(true)},
            msg:      "test",
            msgLevel: level.Warn,
            want:     color.Yellow.Colorize("test"),
        },
        {
            name:     "WithColorization(false)",
            options:  []LoggerOption{WithColorization(false)},
            msg:      "test",
            msgLevel: level.Warn,
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
        msgLevel level.Level
        want     string
    }{
        {
            name: "WithLevelColors()",
            options: []LoggerOption{WithLevelColors(map[level.Level]color.Color{
                level.Debug: color.Yellow,
                level.Info:  color.Blue,
                level.Warn:  color.Green,
                level.Error: color.Red,
                level.Panic: color.Magenta,
            }),
                WithColorization(true)},
            msg:      "test",
            msgLevel: level.Warn,
            want:     color.Green.Colorize("test"),
        },
        {
            name: "With partial level colors",
            options: []LoggerOption{WithLevelColors(map[level.Level]color.Color{
                level.Debug: color.Yellow,
                level.Info:  color.Blue,
                level.Warn:  color.Green,
            }),
                WithColorization(true)},
            msg:      "test",
            msgLevel: level.Warn,
            want:     color.Green.Colorize("test"),
        },
        {
            name: "With partial level colors gets default",
            options: []LoggerOption{WithLevelColors(map[level.Level]color.Color{
                level.Debug: color.Yellow,
                level.Info:  color.Blue,
                level.Warn:  color.Green,
            }),
                WithColorization(true)},
            msg:      "test",
            msgLevel: level.Error,
            want:     color.Red.Colorize("test"),
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
        msgLevel level.Level
        want     string
    }{
        {
            name:     "WithMinLevel(level.Debug), debug shows up",
            options:  []LoggerOption{WithMinLevel(level.Debug)},
            msg:      "test",
            msgLevel: level.Debug,
            want:     "test",
        },
        {
            name:     "WithMinLevel(level.Info), debug doesn't show up",
            options:  []LoggerOption{WithMinLevel(level.Info)},
            msg:      "test",
            msgLevel: level.Debug,
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
        msgLevel level.Level
        want     string
    }{
        {
            name:     "WithPrefixFields()",
            options:  []LoggerOption{WithPrefixFields(field.NewLevelField(bracket.Angle))},
            msg:      "test",
            msgLevel: level.Warn,
            want:     "<WARN> test",
        },
        {
            name:     "Multiple Prefix Fields",
            options:  []LoggerOption{WithPrefixFields(field.NewLevelField(bracket.Angle), field.NewTagField("tag"))},
            msg:      "test",
            msgLevel: level.Warn,
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
        msgLevel level.Level
        want     string
    }{
        {
            name:     "WithSilent(true), error doesn't show up",
            options:  []LoggerOption{WithSilent(true)},
            msg:      "test",
            msgLevel: level.Error,
            want:     "",
        },
        {
            name:     "WithSilent(false), error shows up",
            options:  []LoggerOption{WithSilent(false)},
            msg:      "test",
            msgLevel: level.Error,
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
        msgLevel level.Level
        want     string
    }{
        {
            name:     "WithPrefixFields()",
            options:  []LoggerOption{WithSuffixFields(field.NewLevelField(bracket.Angle))},
            msg:      "test",
            msgLevel: level.Warn,
            want:     "test <WARN>",
        },
        {
            name:     "Multiple Prefix Fields",
            options:  []LoggerOption{WithSuffixFields(field.NewLevelField(bracket.Angle), field.NewTagField("tag"))},
            msg:      "test",
            msgLevel: level.Warn,
            want:     "test <WARN> [tag]",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            testFormatOption(t, tt.options, tt.msg, tt.msgLevel, tt.want)
        })
    }
}
