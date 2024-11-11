package ultralogger

import (
    "bytes"
    "errors"
    "fmt"
    "testing"
)

type invalidField struct{}

func (f invalidField) FieldFormatter() (FieldFormatter, error) {
    return nil, errors.New("invalid field")
}

func Test_ultraFormatter_Format(t *testing.T) {
    type args struct {
        level Level
        msg   string
    }
    tests := []struct {
        name        string
        fields      []Field
        enableColor bool
        levelColors map[Level]Color
        args        args
        want        []byte
        wantErr     bool
    }{
        {
            name: "Default",
            args: args{
                level: Info,
                msg:   "test",
            },
            want: []byte("[tag] <INFO> test"),
            fields: []Field{
                NewTagField(),
                NewLevelField(BracketAngle),
                &FieldMessage{},
            },
        },
        {
            name: "No Fields",
            args: args{
                level: Info,
                msg:   "test",
            },
            want: []byte(""),
        },
        {
            name: "Invalid prefix field throws error",
            args: args{
                level: Info,
                msg:   "test",
            },
            fields: []Field{
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
            fields: []Field{
                &FieldMessage{},
            },
            enableColor: true,
            levelColors: map[Level]Color{
                Debug: ColorWhite,
                Info:  ColorGreen,
                Warn:  ColorYellow,
                Error: ColorRed,
                Panic: ColorMagenta,
            },
            want: ColorGreen.Colorize([]byte("test")),
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
            fields: []Field{
                NewTagField(),
                NewLevelField(BracketAngle),
                &FieldMessage{},
                NewTagField(),
                NewLevelField(BracketAngle),
            },
            want: ColorRed.Colorize([]byte("[tag] <ERROR> test [tag] <ERROR>")),
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            var f LogLineFormatter
            f = &TextFormatter{
                Fields:         tt.fields,
                FieldSeparator: " ",
            }

            if tt.enableColor {
                f = NewColorizedFormatter(f, tt.levelColors)
            }

            messageContext := LogLineContext{
                Level: tt.args.level,
                Tag:   "tag",
            }

            if got, _ := f.FormatLogLine(messageContext, tt.args.msg); !bytes.Equal(got, tt.want) {
                fmt.Println("Got:  ", string(got))
                fmt.Println("Got:  ", got)
                fmt.Println("Want: ", tt.want)
                t.Errorf("Format() = %v, want %v", got, tt.want)
            }
        })
    }
}

//func Test_ultraFormatter_Formatf(t *testing.T) {
//    type args struct {
//        level  Level
//        format string
//        args   []any
//    }
//    tests := []struct {
//        name         string
//        prefixFields []Field
//        suffixFields []Field
//        args         args
//        want         string
//        wantErr      bool
//    }{
//        {
//            name: "Default",
//            args: args{
//                level:  Info,
//                format: "%v %v",
//                args:   []any{"test", "test"},
//            },
//            prefixFields: []Field{
//                NewTagField(),
//                NewLevelField(BracketAngle),
//            },
//            want: "[tag] <INFO> test test",
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            f, err := NewColorizedFormatter(tt.prefixFields, tt.suffixFields, false)
//            if (err != nil) != tt.wantErr {
//                t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantErr)
//                return
//            }
//
//            if err != nil {
//                return
//            }
//
//            if got := f.Formatf(tt.args.level, tt.args.format, tt.args.args...); got != tt.want {
//                t.Errorf("Format() = %v, want %v", got, tt.want)
//            }
//        })
//    }
//}
//
//func Test_ultraFormatter_SetFields(t *testing.T) {
//    tests := []struct {
//        name                    string
//        formatter               Formatter
//        initPrefixFields        []Field
//        initSuffixFields        []Field
//        newPrefixFields         []Field
//        newSuffixFields         []Field
//        msg                     string
//        logLevel                Level
//        before                  string
//        want                    string
//        wantSetPrefixFieldError bool
//        wantSetSuffixFieldError bool
//    }{
//        {
//            name: "SetPrefixFields",
//            initPrefixFields: []Field{
//                NewTagField(),
//                NewLevelField(BracketAngle),
//            },
//            initSuffixFields: []Field{},
//            newPrefixFields: []Field{
//                NewTagField("newTag", WithBracket(BracketRound)),
//            },
//            msg:      "test",
//            logLevel: Info,
//            before:   "[tag] <INFO> test",
//            want:     "(newTag) test",
//        },
//        {
//            name:             "SetSuffixFields",
//            initPrefixFields: []Field{},
//            initSuffixFields: []Field{
//                NewTagField(),
//                NewLevelField(BracketAngle),
//            },
//            newPrefixFields: []Field{},
//            newSuffixFields: []Field{
//                NewTagField("newTag", WithBracket(BracketRound)),
//            },
//            msg:      "test",
//            logLevel: Info,
//            before:   "test [tag] <INFO>",
//            want:     "test (newTag)",
//        },
//        {
//            name:             "SetPrefixFields with existing suffix fields",
//            initPrefixFields: []Field{},
//            initSuffixFields: []Field{
//                NewTagField(),
//                NewLevelField(BracketAngle),
//            },
//            newPrefixFields: []Field{
//                NewTagField("newTag", WithBracket(BracketRound)),
//            },
//            newSuffixFields: nil,
//            msg:             "test",
//            logLevel:        Info,
//            before:          "test [tag] <INFO>",
//            want:            "(newTag) test [tag] <INFO>",
//        },
//        {
//            name: "SetSuffixFields with existing prefix fields",
//            initPrefixFields: []Field{
//                NewTagField(),
//                NewLevelField(BracketAngle),
//            },
//            initSuffixFields: []Field{},
//            newPrefixFields:  nil,
//            newSuffixFields: []Field{
//                NewTagField("newTag", WithBracket(BracketRound)),
//            },
//            msg:      "test",
//            logLevel: Info,
//            before:   "[tag] <INFO> test",
//            want:     "[tag] <INFO> test (newTag)",
//        },
//        {
//            name: "SetPrefixField with invalid field throws error",
//            initPrefixFields: []Field{
//                NewTagField(),
//                NewLevelField(BracketAngle),
//            },
//            initSuffixFields: []Field{},
//            newPrefixFields: []Field{
//                invalidField{},
//            },
//            msg:                     "test",
//            logLevel:                Info,
//            before:                  "[tag] <INFO> test",
//            want:                    "[tag] <INFO> test",
//            wantSetPrefixFieldError: true,
//        },
//        {
//            name:             "SetSuffixField with invalid field throws error",
//            initPrefixFields: []Field{},
//            initSuffixFields: []Field{
//                NewTagField(),
//                NewLevelField(BracketAngle),
//            },
//            newSuffixFields: []Field{
//                invalidField{},
//            },
//            msg:                     "test",
//            logLevel:                Info,
//            before:                  "test [tag] <INFO>",
//            wantSetSuffixFieldError: true,
//        },
//    }
//    for _, tt := range tests {
//        t.Run(tt.name, func(t *testing.T) {
//            f, _ := NewColorizedFormatter(tt.initPrefixFields, tt.initSuffixFields, false)
//
//            beforeLog := f.Format(tt.logLevel, tt.msg)
//            if beforeLog != tt.before {
//                t.Errorf("Format() before log = %v, want %v", beforeLog, tt.before)
//            }
//
//            var err error
//
//            if tt.newPrefixFields != nil {
//                err = f.SetPrefixFields(tt.newPrefixFields...)
//            }
//
//            if (err != nil) != tt.wantSetPrefixFieldError {
//                t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantSetPrefixFieldError)
//                return
//            }
//
//            if err != nil {
//                return
//            }
//
//            if tt.newSuffixFields != nil {
//                err = f.SetSuffixFields(tt.newSuffixFields...)
//            }
//
//            if (err != nil) != tt.wantSetSuffixFieldError {
//                t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantSetSuffixFieldError)
//                return
//            }
//
//            if err != nil {
//                return
//            }
//
//            afterLog := f.Format(tt.logLevel, tt.msg)
//            if afterLog != tt.want {
//                t.Errorf("Format() after log = %v, want %v", afterLog, tt.want)
//            }
//        })
//    }
//}
