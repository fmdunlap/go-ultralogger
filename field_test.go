package ultralogger

import (
    "bytes"
    "fmt"
    "testing"
    "time"
)

type mockClock struct{}

func (c mockClock) Now() time.Time {
    return time.Date(2024, time.November, 7, 19, 30, 0, 0, time.UTC)
}

func TestLevelField(t *testing.T) {
    tests := []struct {
        name       string
        levelField Field
        args       LogLineArgs
        want       string
        wantErr    bool
    }{
        {
            name:       "Default",
            levelField: NewLevelField(BracketAngle),
            args: LogLineArgs{
                Level: Info,
            },
            want: "<INFO>",
        },
        {
            name:       "Round Bracket",
            levelField: NewLevelField(BracketRound),
            args: LogLineArgs{
                Level: Info,
            },
            want: "(INFO)",
        },
        {
            name:       "Debug",
            levelField: NewLevelField(BracketAngle),
            args: LogLineArgs{
                Level: Debug,
            },
            want: "<DEBUG>",
        },
        {
            name:       "Warn",
            levelField: NewLevelField(BracketAngle),
            args: LogLineArgs{
                Level: Warn,
            },
            want: "<WARN>",
        },
        {
            name:       "Error",
            levelField: NewLevelField(BracketAngle),
            args: LogLineArgs{
                Level: Error,
            },
            want: "<ERROR>",
        },
        {
            name:       "Panic",
            levelField: NewLevelField(BracketAngle),
            args: LogLineArgs{
                Level: Panic,
            },
            want: "<PANIC>",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            formatter, err := tt.levelField.FieldFormatter()
            if (err != nil) != tt.wantErr {
                t.Errorf("FieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            result, err := formatter(tt.args, nil)
            if (err != nil) != tt.wantErr {
                t.Errorf("FieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if result.Data != tt.want {
                t.Errorf("FieldFormatter() formatter = %v, want %v", result.Data, tt.want)
            }
        })
    }
}

func TestDateTimeField(t *testing.T) {
    tests := []struct {
        name          string
        dateTimeField *CurrentTimeField
        args          LogLineArgs
        want          string
        wantErr       bool
    }{
        {
            name: "Default",
            dateTimeField: &CurrentTimeField{
                fmtString: "2006-01-02 15:04:05",
                clock:     mockClock{},
            },
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "2024-11-07 19:30:00",
        },
        {
            name: "Only Time",
            dateTimeField: &CurrentTimeField{
                fmtString: "15:04:05",
                clock:     mockClock{},
            },
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "19:30:00",
        },
        {
            name: "Only Date",
            dateTimeField: &CurrentTimeField{
                fmtString: "2006-01-02",
                clock:     mockClock{},
            },
            args: LogLineArgs{
                Level:        Info,
                OutputFormat: OutputFormatText,
            },
            want: "2024-11-07",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            formatter, err := tt.dateTimeField.FieldFormatter()
            if (err != nil) != tt.wantErr {
                t.Errorf("FieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            result, err := formatter(tt.args, nil)
            if (err != nil) != tt.wantErr {
                t.Errorf("FieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if result.Data != tt.want {
                t.Errorf("formatter() got = %v, want %v", result.Data, tt.want)
            }
        })
    }
}

type ComplexMapKey struct {
    Key string
    B   bool
    I   int
}

type ComplexMapValue struct {
    Val string
    B   bool
    I   int
}

func Test_QuickTest(t *testing.T) {
    type tStruct struct {
        Val string
        B   bool
    }

    stringArrayField, _ := NewArrayField[tStruct](
        "stringArray",
        func(args LogLineArgs, data tStruct) any {
            if args.OutputFormat == OutputFormatText {
                return fmt.Sprintf("%s&@&%v", data.Val, data.B)
            }
            return data
        },
    )

    stringField, _ := NewStringField("string")

    boolField, _ := NewBoolField("bool")

    currentTimeField, _ := NewCurrentTimeField("CurrentTime", "2006-01-02 15:04:05")

    responseField, _ := NewResponseField("response", ResponseLogSettings{
        LogStatus: true,
        LogPath:   true,
    })

    mapField, _ := NewMapField[string, string]("map", func(args LogLineArgs, data string) any {
        return data
    }, func(args LogLineArgs, data string) any {
        return data
    })

    complexMapField, _ := NewMapField[ComplexMapKey, ComplexMapValue]("complexMap",
        func(args LogLineArgs, data ComplexMapKey) any {
            if args.OutputFormat != OutputFormatText {
                return fmt.Sprintf("%v:%v", data.Key, data.I)
            }
            return data
        },
        func(args LogLineArgs, data ComplexMapValue) any {
            return data
        },
    )

    testColors := map[Level]Color{
        Debug: ColorAnsiRGB(235, 216, 52),
        Info:  ColorAnsiRGB(12, 240, 228),
        Warn:  ColorAnsiRGB(237, 123, 0),
        Error: ColorAnsiRGB(237, 0, 0),
        Panic: ColorAnsiRGB(237, 0, 0),
    }

    testFormatter, _ := NewFormatter(OutputFormatText, []Field{stringArrayField, stringField, boolField, currentTimeField, mapField, responseField, complexMapField}, WithColorization(testColors))

    buf := &bytes.Buffer{}

    logger, err := NewLoggerWithOptions(WithDestination(buf, testFormatter), WithMinLevel(Debug))
    if err != nil {
        panic(err)
    }

    t.Run("Test", func(t *testing.T) {
        complexMap := map[ComplexMapKey]ComplexMapValue{
            {Key: "testAlpha", B: true, I: 10}: {Val: "ValAlpha", B: true, I: 1},
            {Key: "testBeta", B: false, I: 20}: {Val: "ValBeta", B: false, I: 2},
        }

        logger.Debug(complexMap)
        logger.Info(complexMap)
        logger.Warn(complexMap)
        logger.Error(complexMap)

        time.Sleep(time.Millisecond * 100)

        fmt.Println(buf.String())
    })
}
