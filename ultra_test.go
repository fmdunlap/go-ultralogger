package ultralogger

import (
    "fmt"
    "io"
    "reflect"
    "testing"
    "time"
)

type mockWriter struct {
    buf               []byte
    throwErrorOnWrite bool
}

func (mw *mockWriter) Write(p []byte) (n int, err error) {
    mw.buf = p
    return len(p), nil
}

var staticDate = time.Date(2024, time.November, 7, 19, 30, 0, 0, time.UTC)

type mockClock struct{}

func (c *mockClock) Now() time.Time {
    return staticDate
}

func TestNewUltraLogger(t *testing.T) {
    defaultUltraLogger := &UltraLogger{
        minLevel:         InfoLevel,
        levelBracketType: defaultLevelbracketType,

        tag:            "",
        tagPadSize:     defaultTagSpaceLen,
        padTag:         true,
        tagBracketType: defaultTagBracketType,

        showDate:          true,
        dateFormat:        defaultDateFormat,
        showTime:          true,
        timeFormat:        defaultTimeFormat,
        dateTimeSeparator: defaultDateTimeSeparator,

        colorize:    false,
        levelColors: defaultLevelColors,

        silent: false,

        fallback:               true,
        panicOnPanicLevel:      false,
        computedDateTimeFormat: "2006-01-02 15:04:05",
        computedFmtString:      "%s %s %s",

        clock:  &realClock{},
        writer: io.Discard,
    }

    tests := []struct {
        name   string
        writer io.Writer
        want   *UltraLogger
    }{
        {"Default", io.Discard, defaultUltraLogger},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := NewUltraLogger(tt.writer)
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("NewUltraLogger() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestUltraLogger_OutputWhenMinLevelIs(t *testing.T) {
    tests := []struct {
        name        string
        minLevel    Level
        msg         string
        wantByLevel map[Level]string
    }{
        {
            "Debug",
            DebugLevel,
            "test",
            map[Level]string{
                DebugLevel: "<DEBUG> test\n",
                InfoLevel:  "<INFO> test\n",
                WarnLevel:  "<WARN> test\n",
                ErrorLevel: "<ERROR> test\n",
                PanicLevel: "<PANIC> test\n",
            },
        },
        {
            "Info",
            InfoLevel,
            "test",
            map[Level]string{
                DebugLevel: "",
                InfoLevel:  "<INFO> test\n",
                WarnLevel:  "<WARN> test\n",
                ErrorLevel: "<ERROR> test\n",
                PanicLevel: "<PANIC> test\n",
            },
        },
        {
            "Warn",
            WarnLevel,
            "test",
            map[Level]string{
                DebugLevel: "",
                InfoLevel:  "",
                WarnLevel:  "<WARN> test\n",
                ErrorLevel: "<ERROR> test\n",
                PanicLevel: "<PANIC> test\n",
            },
        },
        {
            "Error",
            ErrorLevel,
            "test",
            map[Level]string{
                DebugLevel: "",
                InfoLevel:  "",
                WarnLevel:  "",
                ErrorLevel: "<ERROR> test\n",
                PanicLevel: "<PANIC> test\n",
            },
        },
        {
            "Panic",
            PanicLevel,
            "test",
            map[Level]string{
                DebugLevel: "",
                InfoLevel:  "",
                WarnLevel:  "",
                ErrorLevel: "",
                PanicLevel: "<PANIC> test\n",
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            writer := &mockWriter{}
            l := &UltraLogger{
                writer:           writer,
                minLevel:         tt.minLevel,
                levelBracketType: BracketTypeAngle,
            }
            l.updateFormatStrings()

            for level, want := range tt.wantByLevel {
                switch level {

                case DebugLevel:
                    l.Debug(tt.msg)
                case InfoLevel:
                    l.Info(tt.msg)
                case WarnLevel:
                    l.Warn(tt.msg)
                case ErrorLevel:
                    l.Error(tt.msg)
                case PanicLevel:
                    l.Panic(tt.msg)
                }

                if got := string(writer.buf); got != want {
                    t.Errorf("%v(%v) => %v, want %v", level, tt.msg, got, want)
                }

                writer.buf = nil
            }
        })
    }
}

func TestUltraLogger_MinLevel(t *testing.T) {
    tests := []struct {
        name  string
        level Level
    }{
        {"MinLevel", InfoLevel},
        {"InvalidMinLevel", Level(42)}, // This will technically pass, but you'll never get log output.
        {"ErrorLevel", ErrorLevel},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockWriter := &mockWriter{}
            l := NewUltraLogger(mockWriter)

            l.SetMinLevel(tt.level)
            if got := l.GetMinLevel(); got != tt.level {
                t.Errorf("UltraLogger.GetMinLevel() = %v, want %v", got, tt.level)
            }
        })
    }
}

func TestUltraLogger_FormattedOutputWhenMinLevelIs(t *testing.T) {
    tests := []struct {
        name        string
        minLevel    Level
        format      string
        args        []any
        wantByLevel map[Level]string
    }{
        {
            "Debug",
            DebugLevel,
            "%s %s",
            []any{"test", "test"},
            map[Level]string{
                DebugLevel: "<DEBUG> test test\n",
                InfoLevel:  "<INFO> test test\n",
                WarnLevel:  "<WARN> test test\n",
                ErrorLevel: "<ERROR> test test\n",
                PanicLevel: "<PANIC> test test\n",
            },
        },
        {
            "Info",
            InfoLevel,
            "%s %s",
            []any{"test", "test"},
            map[Level]string{
                DebugLevel: "",
                InfoLevel:  "<INFO> test test\n",
                WarnLevel:  "<WARN> test test\n",
                ErrorLevel: "<ERROR> test test\n",
                PanicLevel: "<PANIC> test test\n",
            },
        },
        {
            "Warn",
            WarnLevel,
            "%s %s",
            []any{"test", "test"},
            map[Level]string{
                DebugLevel: "",
                InfoLevel:  "",
                WarnLevel:  "<WARN> test test\n",
                ErrorLevel: "<ERROR> test test\n",
                PanicLevel: "<PANIC> test test\n",
            },
        },
        {
            "Error",
            ErrorLevel,
            "%s %s",
            []any{"test", "test"},
            map[Level]string{
                DebugLevel: "",
                InfoLevel:  "",
                WarnLevel:  "",
                ErrorLevel: "<ERROR> test test\n",
                PanicLevel: "<PANIC> test test\n",
            },
        },
        {
            "Panic",
            PanicLevel,
            "%s %s",
            []any{"test", "test"},
            map[Level]string{
                DebugLevel: "",
                InfoLevel:  "",
                WarnLevel:  "",
                ErrorLevel: "",
                PanicLevel: "<PANIC> test test\n",
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            writer := &mockWriter{}
            l := &UltraLogger{
                writer:           writer,
                minLevel:         tt.minLevel,
                levelBracketType: BracketTypeAngle,
            }
            l.updateFormatStrings()

            for level, want := range tt.wantByLevel {
                switch level {

                case DebugLevel:
                    l.Debugf(tt.format, tt.args...)
                case InfoLevel:
                    l.Infof(tt.format, tt.args...)
                case WarnLevel:
                    l.Warnf(tt.format, tt.args...)
                case ErrorLevel:
                    l.Errorf(tt.format, tt.args...)
                case PanicLevel:
                    l.Panicf(tt.format, tt.args...)
                }

                if got := string(writer.buf); got != want {
                    t.Errorf("%v(%v) => %v, want %v", level, fmt.Sprintf(tt.format, tt.args...), got, want)
                }

                writer.buf = nil
            }
        })
    }
}

func TestUltraLogger_LevelBracketType(t *testing.T) {
    tests := []struct {
        name        string
        level       Level
        bracketType BracketType
        msg         string
        want        string
        wantErr     bool
    }{
        {
            name:        "DefaultBracketType",
            level:       InfoLevel,
            bracketType: BracketTypeAngle,
            msg:         "test",
            want:        "<INFO> test",
        },
        {
            name:        "RoundBracketType",
            level:       InfoLevel,
            bracketType: BracketTypeRound,
            msg:         "test",
            want:        "(INFO) test",
        },
        {
            name:        "SquareBracketType",
            level:       InfoLevel,
            bracketType: BracketTypeSquare,
            msg:         "test",
            want:        "[INFO] test",
        },
        {
            name:        "CurlyBracketType",
            level:       InfoLevel,
            bracketType: BracketTypeCurly,
            msg:         "test",
            want:        "{INFO} test",
        },
        {
            name:        "UnknownBracketType",
            level:       InfoLevel,
            bracketType: BracketType(42),
            msg:         "test",
            want:        "<INFO> test",
            wantErr:     true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockWriter := &mockWriter{}
            l := NewUltraLogger(mockWriter)
            l.ShowDate(false)
            l.ShowTime(false)

            _, err := l.SetLevelBracketType(tt.bracketType)
            if err != nil {
                if !tt.wantErr {
                    t.Errorf("UltraLogger.SetLevelBracketType() error = %v, wantErr %v", err, tt.wantErr)
                }
                return
            }

            if got := l.GetLevelBracketType(); got != tt.bracketType {
                t.Errorf("UltraLogger.GetLevelBracketType() = %v, want %v", got, tt.bracketType)
            }

            if got := l.Slog(tt.level, tt.msg); got != tt.want {
                t.Errorf("UltraLogger.Slog() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestUltraLogger_Tag(t *testing.T) {
    tests := []struct {
        name              string
        tag               string
        msg               string
        padTag            bool
        tagPadSize        int
        tagBracketType    BracketType
        showDate          bool
        showTime          bool
        want              string
        wantTagBracketErr bool
    }{
        {
            name:           "DefaultTag, no padding",
            tag:            "test",
            tagBracketType: BracketTypeSquare,
            msg:            "test",
            want:           "[test] <INFO> test",
        },
        {
            name:           "Default bracket, padding",
            tag:            "test",
            tagBracketType: BracketTypeSquare,
            msg:            "test",
            padTag:         true,
            tagPadSize:     10,
            want:           "[test]    <INFO> test",
        },
        {
            name:           "Default bracket, no padding, with date time",
            tag:            "test",
            tagBracketType: BracketTypeSquare,
            msg:            "test",
            showDate:       true,
            showTime:       true,
            want:           "2024-11-07 19:30:00 [test] <INFO> test",
        },
        {
            name:           "Default bracket, padding, with date time",
            tag:            "test",
            tagBracketType: BracketTypeSquare,
            msg:            "test",
            showDate:       true,
            showTime:       true,
            padTag:         true,
            tagPadSize:     10,
            want:           "2024-11-07 19:30:00 [test]    <INFO> test",
        },
        {
            name:       "Default bracket, padding, with date, no tag",
            tag:        "",
            msg:        "test",
            showDate:   true,
            padTag:     true,
            tagPadSize: 10,
            want:       "2024-11-07 <INFO> test",
        },
        {
            name:           "None bracket, no padding, no date time",
            tag:            "test",
            tagBracketType: BracketTypeNone,
            msg:            "test",
            want:           "test <INFO> test",
        },
        {
            name:           "None bracket, padding, no date time",
            tag:            "test",
            tagBracketType: BracketTypeNone,
            msg:            "test",
            padTag:         true,
            tagPadSize:     10,
            want:           "test      <INFO> test",
        },
        {
            name:           "None bracket, no padding, with date time",
            tag:            "test",
            tagBracketType: BracketTypeNone,
            msg:            "test",
            showDate:       true,
            showTime:       true,
            want:           "2024-11-07 19:30:00 test <INFO> test",
        },
        {
            name:           "None bracket, padding, with date time",
            tag:            "test",
            tagBracketType: BracketTypeNone,
            msg:            "test",
            showDate:       true,
            showTime:       true,
            padTag:         true,
            tagPadSize:     10,
            want:           "2024-11-07 19:30:00 test      <INFO> test",
        },
        {
            name:           "Different tag bracket type",
            tag:            "test",
            tagBracketType: BracketTypeRound,
            msg:            "test",
            want:           "(test) <INFO> test",
        },
        {
            name:              "Invalid tag bracket type throws error",
            tag:               "test",
            tagBracketType:    BracketType(42),
            msg:               "test",
            wantTagBracketErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockWriter := &mockWriter{}

            l := &UltraLogger{
                writer:            mockWriter,
                minLevel:          DebugLevel,
                levelBracketType:  BracketTypeAngle,
                tag:               "",
                showDate:          tt.showDate,
                dateFormat:        defaultDateFormat,
                showTime:          tt.showTime,
                timeFormat:        defaultTimeFormat,
                dateTimeSeparator: defaultDateTimeSeparator,

                clock: &mockClock{},
            }
            l.updateFormatStrings()

            l.SetTag(tt.tag)

            if tt.padTag {
                l.EnabledTagPadding(true)
                l.SetTagPadSize(tt.tagPadSize)
            }

            _, err := l.SetTagBracketType(tt.tagBracketType)
            if err != nil {
                if !tt.wantTagBracketErr {
                    t.Errorf("UltraLogger.SetTagBracketType() error = %v, wantErr %v", err, tt.wantTagBracketErr)
                    return
                }
                return
            }

            if tt.tagBracketType != l.GetTagBracketType() {
                t.Errorf("UltraLogger.GetTagBracketType() = %v, want %v", l.GetTagBracketType(), tt.tagBracketType)
            }

            if got := l.GetTag(); got != tt.tag {
                t.Errorf("UltraLogger.GetTag() = %v, want %v", got, tt.tag)
            }

            if got := l.Slog(InfoLevel, tt.msg); got != tt.want {
                t.Errorf("UltraLogger.Slog() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestUltraLogger_TimeDate(t *testing.T) {
    tagBracketType := BracketTypeSquare
    levelBracketType := BracketTypeAngle

    tests := []struct {
        name string
        msg  string

        tag          string
        enablePadTag bool
        tagPadSize   int

        customDateFormat        string
        disableDate             bool
        customTimeFormat        string
        disableTime             bool
        customDateTimeSeparator string

        logLevel Level

        want          string
        wantDateError bool
        wantTimeError bool
    }{
        {
            name:     "Standard format, no tag, debug",
            msg:      "test",
            tag:      "",
            logLevel: DebugLevel,
            want:     "2024-11-07 19:30:00 <DEBUG> test\n",
        },
        {
            name:     "Standard format, no tag, info",
            msg:      "test",
            tag:      "",
            logLevel: InfoLevel,
            want:     "2024-11-07 19:30:00 <INFO> test\n",
        },
        {
            name:     "Standard format, tag, debug",
            msg:      "test",
            tag:      "test",
            logLevel: DebugLevel,
            want:     "2024-11-07 19:30:00 [test] <DEBUG> test\n",
        },
        {
            name:         "Standard format, tag padding, info",
            msg:          "test",
            tag:          "test",
            logLevel:     InfoLevel,
            enablePadTag: true,
            tagPadSize:   10,
            want:         "2024-11-07 19:30:00 [test]    <INFO> test\n",
        },
        {
            name:        "Standard format, no date, tag, info",
            msg:         "test",
            tag:         "test",
            logLevel:    InfoLevel,
            disableDate: true,
            want:        "19:30:00 [test] <INFO> test\n",
        },
        {
            name:        "Standard format, no time, tag, info",
            msg:         "test",
            tag:         "test",
            logLevel:    InfoLevel,
            disableTime: true,
            want:        "2024-11-07 [test] <INFO> test\n",
        },
        {
            name:        "Standard format, no date and time, tag, info",
            msg:         "test",
            tag:         "test",
            logLevel:    InfoLevel,
            disableDate: true,
            disableTime: true,
            want:        "[test] <INFO> test\n",
        },
        {
            name:         "Standard format, no date and time, tag padding, info",
            msg:          "test",
            tag:          "test",
            logLevel:     InfoLevel,
            disableDate:  true,
            disableTime:  true,
            enablePadTag: true,
            tagPadSize:   10,
            want:         "[test]    <INFO> test\n",
        },
        {
            name:                    "Standard format, no tag, info, with custom date time separator",
            msg:                     "test",
            tag:                     "",
            logLevel:                InfoLevel,
            customDateTimeSeparator: "@",
            want:                    "2024-11-07@19:30:00 <INFO> test\n",
        },
        {
            name:                    "Standard format, tag, info, no date, with custom date time separator",
            msg:                     "test",
            tag:                     "test",
            logLevel:                InfoLevel,
            disableDate:             true,
            customDateTimeSeparator: "@",
            want:                    "19:30:00 [test] <INFO> test\n",
        },
        {
            name:                    "Standard format, tag, info, no time, with custom date time separator",
            msg:                     "test",
            tag:                     "test",
            logLevel:                InfoLevel,
            disableTime:             true,
            customDateTimeSeparator: "@",
            want:                    "2024-11-07 [test] <INFO> test\n",
        },
        {
            name:                    "Standard format, tag, info, no date and time, with custom date time separator",
            msg:                     "test",
            tag:                     "test",
            logLevel:                InfoLevel,
            disableDate:             true,
            disableTime:             true,
            customDateTimeSeparator: "@",
            want:                    "[test] <INFO> test\n",
        },
        {
            name:             "Modified date format, no tag, info",
            msg:              "test",
            tag:              "",
            logLevel:         InfoLevel,
            customDateFormat: "01-02-06",
            want:             "11-07-24 19:30:00 <INFO> test\n",
        },
        {
            name:             "Modified date format, tag, info",
            msg:              "test",
            tag:              "test",
            logLevel:         InfoLevel,
            customDateFormat: "01-02-06",
            want:             "11-07-24 19:30:00 [test] <INFO> test\n",
        },
        {
            name:             "Modified time format, no tag, info",
            msg:              "test",
            tag:              "",
            logLevel:         InfoLevel,
            customTimeFormat: "15.04.05",
            want:             "2024-11-07 19.30.00 <INFO> test\n",
        },
        {
            name:             "Modified time format, tag, info",
            msg:              "test",
            tag:              "tag",
            logLevel:         InfoLevel,
            customTimeFormat: "15.04.05",
            want:             "2024-11-07 19.30.00 [tag] <INFO> test\n",
        },
        {
            name:             "Modified date and time format, no tag, info",
            msg:              "test",
            tag:              "",
            logLevel:         InfoLevel,
            customDateFormat: "01-02-06",
            customTimeFormat: "15.04.05",
            want:             "11-07-24 19.30.00 <INFO> test\n",
        },
        {
            name:             "Modified date and time format, tag, info",
            msg:              "test",
            tag:              "tag",
            logLevel:         InfoLevel,
            customDateFormat: "01-02-06",
            customTimeFormat: "15.04.05",
            want:             "11-07-24 19.30.00 [tag] <INFO> test\n",
        },
        {
            name:                    "Modified date and time format, no tag, info, with custom date time separator",
            msg:                     "test",
            tag:                     "",
            logLevel:                InfoLevel,
            customDateFormat:        "01-02-06",
            customTimeFormat:        "15.04.05",
            customDateTimeSeparator: "@",
            want:                    "11-07-24@19.30.00 <INFO> test\n",
        },
        {
            name:                    "Modified date format, no time, no tag, info, with custom date time separator",
            msg:                     "test",
            tag:                     "test",
            logLevel:                InfoLevel,
            disableTime:             true,
            customDateFormat:        "01-02-06",
            customDateTimeSeparator: "@",
            want:                    "11-07-24 [test] <INFO> test\n",
        },
        {
            name:                    "Modified time format, no date, tag, info, no time, with custom date time separator",
            msg:                     "test",
            tag:                     "test",
            logLevel:                InfoLevel,
            disableDate:             true,
            customTimeFormat:        "15.04.05",
            customDateTimeSeparator: "@",
            want:                    "19.30.00 [test] <INFO> test\n",
        },
        {
            name:             "Invalid date format throws error",
            msg:              "test",
            logLevel:         InfoLevel,
            customDateFormat: "blah",
            wantDateError:    true,
        },
        {
            name:             "Invalid time format throws error",
            msg:              "test",
            logLevel:         InfoLevel,
            customTimeFormat: "blah",
            wantTimeError:    true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            writer := &mockWriter{}
            l := &UltraLogger{
                writer:            writer,
                minLevel:          DebugLevel,
                levelBracketType:  levelBracketType,
                tag:               tt.tag,
                tagBracketType:    tagBracketType,
                showDate:          true,
                dateFormat:        defaultDateFormat,
                showTime:          true,
                timeFormat:        defaultTimeFormat,
                dateTimeSeparator: defaultDateTimeSeparator,

                clock: &mockClock{},
            }
            l.updateFormatStrings()

            if tt.enablePadTag {
                l.EnabledTagPadding(true)
                l.SetTagPadSize(tt.tagPadSize)

                if l.GetTagPaddingEnabled() != true {
                    t.Errorf("UltraLogger.EnabledTagPadding() got = %v, want %v", l.GetTagPaddingEnabled(), true)
                }

                if l.GetTagPadSize() != tt.tagPadSize {
                    t.Errorf("UltraLogger.SetTagPadSize() got = %v, want %v", l.GetTagPadSize(), tt.tagPadSize)
                }
            }

            if tt.customDateFormat != "" {
                _, err := l.SetDateFormat(tt.customDateFormat)
                if tt.wantDateError {
                    if err == nil {
                        t.Errorf("UltraLogger.SetDateFormat() error = %v, wantErr %v", err, tt.wantDateError)
                    }
                    return
                }

                if l.GetDateFormat() != tt.customDateFormat {
                    t.Errorf("UltraLogger.SetDateFormat() got = %v, want %v", l.GetDateFormat(), tt.customDateFormat)
                }
            }

            if tt.customTimeFormat != "" {
                _, err := l.SetTimeFormat(tt.customTimeFormat)
                if tt.wantTimeError {
                    if err == nil {
                        t.Errorf("UltraLogger.SetTimeFormat() error = %v, wantErr %v", err, tt.wantTimeError)
                    }
                    return
                }

                if l.GetTimeFormat() != tt.customTimeFormat {
                    t.Errorf("UltraLogger.SetTimeFormat() got = %v, want %v", l.GetTimeFormat(), tt.customTimeFormat)
                }
            }

            if tt.customDateTimeSeparator != "" {
                l.SetDateTimeSeparator(tt.customDateTimeSeparator)
                if l.GetDateTimeSeparator() != tt.customDateTimeSeparator {
                    t.Errorf("UltraLogger.SetDateTimeSeparator() got = %v, want %v", l.GetDateTimeSeparator(), tt.customDateTimeSeparator)
                }
            }

            if tt.disableDate {
                l.ShowDate(false)

                if l.GetShowDate() != false {
                    t.Errorf("UltraLogger.ShowDate() got = %v, want %v", l.GetShowDate(), true)
                }
            }

            if tt.disableTime {
                l.ShowTime(false)

                if l.GetShowTime() != false {
                    t.Errorf("UltraLogger.ShowTime() got = %v, want %v", l.GetShowTime(), true)
                }
            }

            switch tt.logLevel {
            case DebugLevel:
                l.Debug(tt.msg)
            case InfoLevel:
                l.Info(tt.msg)
            case WarnLevel:
                l.Warn(tt.msg)
            case ErrorLevel:
                l.Error(tt.msg)
            case PanicLevel:
                l.Panic(tt.msg)
            }

            if got := string(writer.buf); got != tt.want {
                t.Errorf("UltraLogger.TimeDate() => %v, want %v", got, tt.want)
            }
        })
    }
}

func TestUltraLogger_ValidDateFormat(t *testing.T) {
    tests := []struct {
        name   string
        format string
        want   bool
    }{
        {
            name:   "Valid date format",
            format: "2006-01-02",
            want:   true,
        },
        {
            name:   "Valid date format, with time",
            format: "2006-01-02 15:04:05",
            want:   true,
        },
        {
            name:   "Invalid date format",
            format: "blah",
            want:   false,
        },
        {
            name:   "Invalid date format, with time",
            format: "blah 15:04:05",
            want:   false,
        },
        // This is fine behavior, but it's not ideal. If yah find a better way to do this, please PR.
        {
            name:   "Valid date format, with invalid time",
            format: "2006-01-02 blah",
            want:   true,
        },
        {
            name:   "Partially invalid date format",
            format: "2024-01-02 15:04:05",
            want:   false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := validDateFormat(tt.format); got != tt.want {
                t.Errorf("validDateFormat() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestUltraLogger_ValidTimeFormat(t *testing.T) {
    tests := []struct {
        name   string
        format string
        want   bool
    }{
        {
            name:   "Valid time format",
            format: "15:04:05",
            want:   true,
        },
        {
            name:   "Invalid time format",
            format: "blah",
            want:   false,
        },
        {
            name:   "Invalid time format, valid date",
            format: "2006-01-02 23:04:05",
            want:   false,
        },
        {
            name:   "Valid time format, valid date",
            format: "2006-01-02 15:04:05",
            want:   true,
        },
        // This is fine behavior, but it's not ideal. If yah find a better way to do this, please PR.
        {
            name:   "Valid time format, invalid date",
            format: "blah 15:04:05",
            want:   true,
        },
        {
            name:   "Partially invalid time format",
            format: "23:04:05",
            want:   false,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := validTimeFormat(tt.format); got != tt.want {
                t.Errorf("validTimeFormat() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestUltraLogger_Silent(t *testing.T) {
    tests := []struct {
        name        string
        want        bool
        msg         string
        wantMessage string
    }{
        {name: "Silent", want: true, msg: "test", wantMessage: ""},
        {name: "Not Silent", want: false, msg: "test", wantMessage: "<WARN> test\n"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockWriter := &mockWriter{}
            l := NewUltraLogger(mockWriter)
            l.ShowDate(false)
            l.ShowTime(false)

            l.SetSilent(tt.want)

            if got := l.GetSilent(); got != tt.want {
                t.Errorf("UltraLogger.GetSilent() = %v, want %v", got, tt.want)
            }

            l.Warn(tt.msg)
            if got := string(mockWriter.buf); got != tt.wantMessage {
                t.Errorf("UltraLogger.Warn() = %v, want %v", got, tt.wantMessage)
            }
        })
    }
}

func TestUltraLogger_PanicOnPanicLevel(t *testing.T) {
    tests := []struct {
        name        string
        want        bool
        msg         string
        wantMessage string
    }{
        {name: "PanicOnPanicLevel", want: true, msg: "test", wantMessage: "<PANIC> test\n"},
        {name: "Not PanicOnPanicLevel", want: false, msg: "test", wantMessage: "<PANIC> test\n"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockWriter := &mockWriter{}
            l := NewUltraLogger(mockWriter)
            l.ShowDate(false)
            l.ShowTime(false)

            l.SetPanicOnPanicLevel(tt.want)

            if l.GetPanicOnPanicLevel() != tt.want {
                t.Errorf("UltraLogger.GetPanicOnPanicLevel() = %v, want %v", l.GetPanicOnPanicLevel(), tt.want)
            }

            defer func() {
                if got := string(mockWriter.buf); got != tt.wantMessage {
                    t.Errorf("UltraLogger.Panic() = %v, want %v", got, tt.wantMessage)
                }
                if didPanic := recover() != nil; didPanic != tt.want {
                    t.Errorf("UltraLogger.Panic() did not panic as expected, got %v, want %v", didPanic, tt.want)
                }
            }()

            l.Panic(tt.msg)
        })
    }
}

func TestUltraLogger_Colors(t *testing.T) {
    tests := []struct {
        name                 string
        enableColor          bool
        msg                  string
        msgLevel             Level
        setCustomColor       bool
        customColor          Color
        wantCustomColorError bool
        wantColor            Color
        wantMsg              string
    }{
        {
            name:        "Default warn",
            enableColor: true,
            msg:         "test",
            msgLevel:    WarnLevel,
            wantColor:   ColorYellow,
        },
        {
            name:        "Default error",
            enableColor: true,
            msg:         "test",
            msgLevel:    ErrorLevel,
            wantColor:   ColorRed,
        },
        {
            name:        "Default panic",
            enableColor: true,
            msg:         "test",
            msgLevel:    PanicLevel,
            wantColor:   ColorMagenta,
        },
        {
            name:           "Custom color",
            enableColor:    true,
            setCustomColor: true,
            customColor:    ColorGreen,
            msg:            "test",
            msgLevel:       WarnLevel,
            wantColor:      ColorGreen,
        },
        {
            name:                 "Custom color, invalid level, throws error",
            enableColor:          true,
            setCustomColor:       true,
            customColor:          ColorGreen,
            msg:                  "test",
            msgLevel:             42,
            wantColor:            ColorGreen,
            wantCustomColorError: true,
        },
        {
            name:                 "Invalid color, throws error",
            enableColor:          true,
            msg:                  "test",
            msgLevel:             WarnLevel,
            setCustomColor:       true,
            customColor:          Color(42),
            wantCustomColorError: true,
            wantMsg:              "<WARN> test\n",
        },
        {
            name:        "Test disabled color",
            enableColor: false,
            msg:         "test",
            msgLevel:    WarnLevel,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockWriter := &mockWriter{}
            l := NewUltraLogger(mockWriter)
            l.ShowDate(false)
            l.ShowTime(false)

            l.SetColorize(tt.enableColor)

            if got := l.GetColorize(); got != tt.enableColor {
                t.Errorf("UltraLogger.GetColorize() = %v, want %v", got, tt.enableColor)
            }

            if tt.setCustomColor {
                _, err := l.SetLevelColor(tt.msgLevel, tt.customColor)
                if tt.wantCustomColorError {
                    if err == nil {
                        t.Errorf("UltraLogger.SetLevelColor() error = %v, wantErr %v", err, tt.wantCustomColorError)
                    }
                    return
                }

                if l.GetLevelColor(tt.msgLevel) != tt.customColor {
                    t.Errorf("UltraLogger.SetLevelColor() got = %v, want %v", l.GetLevelColor(tt.msgLevel), tt.customColor)
                }
            }

            switch tt.msgLevel {
            case DebugLevel:
                l.Debug(tt.msg)
            case InfoLevel:
                l.Info(tt.msg)
            case WarnLevel:
                l.Warn(tt.msg)
            case ErrorLevel:
                l.Error(tt.msg)
            case PanicLevel:
                l.Panic(tt.msg)
            }

            // Bite me. It's the same system that the logger (currently) uses to format the output.
            //
            // If you're hitting this in the future, and you changed how colorizing works, then check if you need to
            // change this.
            expected := testExpectedColorize(tt.wantColor, tt.msgLevel, tt.msg)
            if !tt.enableColor {
                expected = BracketTypeAngle.wrap(tt.msgLevel.String()) + " " + tt.msg + "\n"
            }

            if got := string(mockWriter.buf); got != expected {
                t.Errorf("UltraLogger.Warn() = '%v', want '%v'", got, expected)
            }

        })
    }
}

func TestUltraLogger_MultipleColors(t *testing.T) {
    tests := []struct {
        name           string
        msg            string
        customColors   map[Level]Color
        expectedColors map[Level]Color
        wantErr        bool
    }{
        {
            name: "Custom colors",
            msg:  "test",
            customColors: map[Level]Color{
                DebugLevel: ColorMagenta,
                InfoLevel:  ColorBlue,
                WarnLevel:  ColorGray,
                ErrorLevel: ColorGreen,
                PanicLevel: ColorCyan,
            },
            expectedColors: map[Level]Color{
                DebugLevel: ColorMagenta,
                InfoLevel:  ColorBlue,
                WarnLevel:  ColorGray,
                ErrorLevel: ColorGreen,
                PanicLevel: ColorCyan,
            },
        },
        {
            name: "Partially custom colors",
            msg:  "different",
            customColors: map[Level]Color{
                DebugLevel: ColorMagenta,
                InfoLevel:  ColorBlue,
                WarnLevel:  ColorGray,
            },
            expectedColors: map[Level]Color{
                DebugLevel: ColorMagenta,
                InfoLevel:  ColorBlue,
                WarnLevel:  ColorGray,
                ErrorLevel: ColorRed,
                PanicLevel: ColorMagenta,
            },
        },
        {
            name: "Invalid color",
            msg:  "different",
            customColors: map[Level]Color{
                DebugLevel: ColorMagenta,
                InfoLevel:  ColorBlue,
                WarnLevel:  ColorGray,
                ErrorLevel: Color(42),
                PanicLevel: Color(42),
            },
            wantErr: true,
        },
        {
            name: "Invalid level",
            msg:  "different",
            customColors: map[Level]Color{
                DebugLevel: ColorMagenta,
                InfoLevel:  ColorBlue,
                WarnLevel:  ColorGray,
                ErrorLevel: Color(42),
                PanicLevel: Color(42),
            },
            wantErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockWriter := &mockWriter{}
            l := NewUltraLogger(mockWriter)
            l.SetMinLevel(DebugLevel)
            l.ShowDate(false)
            l.ShowTime(false)

            for level, color := range l.GetLevelColors() {
                if color != defaultLevelColors[level] {
                    t.Errorf("UltraLogger.GetLevelColors() got = %v, want %v", color, defaultLevelColors[level])
                }
            }

            l.SetColorize(true)
            _, err := l.SetLevelColors(tt.customColors)

            for level, color := range tt.customColors {
                fmt.Println(level, color)
            }

            if tt.wantErr {
                if err == nil {
                    t.Errorf("UltraLogger.SetLevelColors() error = %v, wantErr %v", err, tt.wantErr)
                }
                return
            } else if err != nil {
                t.Errorf("UltraLogger.SetLevelColors() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            for _, level := range AllLevels() {
                switch level {
                case DebugLevel:
                    l.Debug(tt.msg)
                case InfoLevel:
                    l.Info(tt.msg)
                case WarnLevel:
                    l.Warn(tt.msg)
                case ErrorLevel:
                    l.Error(tt.msg)
                case PanicLevel:
                    l.Panic(tt.msg)
                }

                expected := testExpectedColorize(tt.expectedColors[level], level, tt.msg)
                if got := string(mockWriter.buf); got != expected {
                    t.Errorf("UltraLogger.Warn() = '%v', want '%v'", got, expected)

                }
            }
        })
    }
}

func testExpectedColorize(color Color, msgLevel Level, msg string) string {
    return colorize(color, BracketTypeAngle.wrap(msgLevel.String())) + " " + colorize(color, msg) + "\n"
}
