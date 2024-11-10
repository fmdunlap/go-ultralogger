package ultralogger

import (
    "errors"
    "testing"
)

type invalidField struct{}

func (f invalidField) FieldPrinter() (FieldPrinterFunc, error) {
    return nil, errors.New("invalid field")
}

func Test_ultraFormatter_Format(t *testing.T) {
    type args struct {
        level Level
        msg   string
    }
    tests := []struct {
        name         string
        prefixFields []Field
        suffixFields []Field
        enableColor  bool
        levelColors  map[Level]Color
        args         args
        want         string
        wantErr      bool
    }{
        {
            name: "Default",
            args: args{
                level: Info,
                msg:   "test",
            },
            want: "[tag] <INFO> test",
            prefixFields: []Field{
                NewTagField("tag"),
                NewLevelField(BracketAngle),
            },
        },
        {
            name: "Suffix Fields",
            args: args{
                level: Info,
                msg:   "test",
            },
            want: "test [tag] <INFO>",
            suffixFields: []Field{
                NewTagField("tag"),
                NewLevelField(BracketAngle),
            },
        },
        {
            name: "No Fields",
            args: args{
                level: Info,
                msg:   "test",
            },
            want: "test",
        },
        {
            name: "Invalid prefix field throws error",
            args: args{
                level: Info,
                msg:   "test",
            },
            prefixFields: []Field{
                invalidField{},
            },
            wantErr: true,
        },
        {
            name: "Invalid suffix field throws error",
            args: args{
                level: Info,
                msg:   "test",
            },
            suffixFields: []Field{
                invalidField{},
            },
            wantErr: true,
        },
        {
            name: "Colorize",
            args: args{
                level: Info,
                msg:   "test",
            },
            enableColor: true,
            levelColors: map[Level]Color{
                Debug: ColorWhite,
                Info:  ColorGreen,
                Warn:  ColorYellow,
                Error: ColorRed,
                Panic: ColorMagenta,
            },
            want: ColorGreen.Colorize("test"),
        },
        {
            name: "Colorize fields",
            args: args{
                level: Error,
                msg:   "test",
            },
            enableColor: true,
            levelColors: map[Level]Color{
                Debug: ColorWhite,
                Info:  ColorGreen,
                Warn:  ColorYellow,
                Error: ColorRed,
                Panic: ColorMagenta,
            },
            prefixFields: []Field{
                NewTagField("tag"),
                NewLevelField(BracketAngle),
            },
            suffixFields: []Field{
                NewTagField("tag"),
                NewLevelField(BracketAngle),
            },
            want: ColorRed.Colorize("[tag] <ERROR> test [tag] <ERROR>"),
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f, err := NewColorizedFormatter(tt.prefixFields, tt.suffixFields, tt.enableColor)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if err != nil {
                return
            }

            _ = f.SetLevelColors(tt.levelColors)

            if got := f.Format(tt.args.level, tt.args.msg); got != tt.want {
                t.Errorf("Format() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_ultraFormatter_Formatf(t *testing.T) {
    type args struct {
        level  Level
        format string
        args   []any
    }
    tests := []struct {
        name         string
        prefixFields []Field
        suffixFields []Field
        args         args
        want         string
        wantErr      bool
    }{
        {
            name: "Default",
            args: args{
                level:  Info,
                format: "%v %v",
                args:   []any{"test", "test"},
            },
            prefixFields: []Field{
                NewTagField("tag"),
                NewLevelField(BracketAngle),
            },
            want: "[tag] <INFO> test test",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f, err := NewColorizedFormatter(tt.prefixFields, tt.suffixFields, false)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if err != nil {
                return
            }

            if got := f.Formatf(tt.args.level, tt.args.format, tt.args.args...); got != tt.want {
                t.Errorf("Format() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_ultraFormatter_SetFields(t *testing.T) {
    tests := []struct {
        name                    string
        formatter               Formatter
        initPrefixFields        []Field
        initSuffixFields        []Field
        newPrefixFields         []Field
        newSuffixFields         []Field
        msg                     string
        logLevel                Level
        before                  string
        want                    string
        wantSetPrefixFieldError bool
        wantSetSuffixFieldError bool
    }{
        {
            name: "SetPrefixFields",
            initPrefixFields: []Field{
                NewTagField("tag"),
                NewLevelField(BracketAngle),
            },
            initSuffixFields: []Field{},
            newPrefixFields: []Field{
                NewTagField("newTag", WithBracket(BracketRound)),
            },
            msg:      "test",
            logLevel: Info,
            before:   "[tag] <INFO> test",
            want:     "(newTag) test",
        },
        {
            name:             "SetSuffixFields",
            initPrefixFields: []Field{},
            initSuffixFields: []Field{
                NewTagField("tag"),
                NewLevelField(BracketAngle),
            },
            newPrefixFields: []Field{},
            newSuffixFields: []Field{
                NewTagField("newTag", WithBracket(BracketRound)),
            },
            msg:      "test",
            logLevel: Info,
            before:   "test [tag] <INFO>",
            want:     "test (newTag)",
        },
        {
            name:             "SetPrefixFields with existing suffix fields",
            initPrefixFields: []Field{},
            initSuffixFields: []Field{
                NewTagField("tag"),
                NewLevelField(BracketAngle),
            },
            newPrefixFields: []Field{
                NewTagField("newTag", WithBracket(BracketRound)),
            },
            newSuffixFields: nil,
            msg:             "test",
            logLevel:        Info,
            before:          "test [tag] <INFO>",
            want:            "(newTag) test [tag] <INFO>",
        },
        {
            name: "SetSuffixFields with existing prefix fields",
            initPrefixFields: []Field{
                NewTagField("tag"),
                NewLevelField(BracketAngle),
            },
            initSuffixFields: []Field{},
            newPrefixFields:  nil,
            newSuffixFields: []Field{
                NewTagField("newTag", WithBracket(BracketRound)),
            },
            msg:      "test",
            logLevel: Info,
            before:   "[tag] <INFO> test",
            want:     "[tag] <INFO> test (newTag)",
        },
        {
            name: "SetPrefixField with invalid field throws error",
            initPrefixFields: []Field{
                NewTagField("tag"),
                NewLevelField(BracketAngle),
            },
            initSuffixFields: []Field{},
            newPrefixFields: []Field{
                invalidField{},
            },
            msg:                     "test",
            logLevel:                Info,
            before:                  "[tag] <INFO> test",
            want:                    "[tag] <INFO> test",
            wantSetPrefixFieldError: true,
        },
        {
            name:             "SetSuffixField with invalid field throws error",
            initPrefixFields: []Field{},
            initSuffixFields: []Field{
                NewTagField("tag"),
                NewLevelField(BracketAngle),
            },
            newSuffixFields: []Field{
                invalidField{},
            },
            msg:                     "test",
            logLevel:                Info,
            before:                  "test [tag] <INFO>",
            wantSetSuffixFieldError: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f, _ := NewColorizedFormatter(tt.initPrefixFields, tt.initSuffixFields, false)

            beforeLog := f.Format(tt.logLevel, tt.msg)
            if beforeLog != tt.before {
                t.Errorf("Format() before log = %v, want %v", beforeLog, tt.before)
            }

            var err error

            if tt.newPrefixFields != nil {
                err = f.SetPrefixFields(tt.newPrefixFields...)
            }

            if (err != nil) != tt.wantSetPrefixFieldError {
                t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantSetPrefixFieldError)
                return
            }

            if err != nil {
                return
            }

            if tt.newSuffixFields != nil {
                err = f.SetSuffixFields(tt.newSuffixFields...)
            }

            if (err != nil) != tt.wantSetSuffixFieldError {
                t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantSetSuffixFieldError)
                return
            }

            if err != nil {
                return
            }

            afterLog := f.Format(tt.logLevel, tt.msg)
            if afterLog != tt.want {
                t.Errorf("Format() after log = %v, want %v", afterLog, tt.want)
            }
        })
    }
}
