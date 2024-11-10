package formatter

import (
    "errors"
    "github.com/fmdunlap/go-ultralogger/bracket"
    "github.com/fmdunlap/go-ultralogger/color"
    "github.com/fmdunlap/go-ultralogger/field"
    "github.com/fmdunlap/go-ultralogger/level"
    "testing"
)

type invalidField struct{}

func (f invalidField) FieldPrinter() (field.FieldPrinterFunc, error) {
    return nil, errors.New("invalid field")
}

func Test_ultraFormatter_Format(t *testing.T) {
    type args struct {
        level level.Level
        msg   string
    }
    tests := []struct {
        name         string
        prefixFields []field.Field
        suffixFields []field.Field
        enableColor  bool
        levelColors  map[level.Level]color.Color
        args         args
        want         string
        wantErr      bool
    }{
        {
            name: "Default",
            args: args{
                level: level.Info,
                msg:   "test",
            },
            want: "[tag] <INFO> test",
            prefixFields: []field.Field{
                field.NewTagField("tag"),
                field.NewLevelField(bracket.Angle),
            },
        },
        {
            name: "Suffix Fields",
            args: args{
                level: level.Info,
                msg:   "test",
            },
            want: "test [tag] <INFO>",
            suffixFields: []field.Field{
                field.NewTagField("tag"),
                field.NewLevelField(bracket.Angle),
            },
        },
        {
            name: "No Fields",
            args: args{
                level: level.Info,
                msg:   "test",
            },
            want: "test",
        },
        {
            name: "Invalid prefix field throws error",
            args: args{
                level: level.Info,
                msg:   "test",
            },
            prefixFields: []field.Field{
                invalidField{},
            },
            wantErr: true,
        },
        {
            name: "Invalid suffix field throws error",
            args: args{
                level: level.Info,
                msg:   "test",
            },
            suffixFields: []field.Field{
                invalidField{},
            },
            wantErr: true,
        },
        {
            name: "Colorize",
            args: args{
                level: level.Info,
                msg:   "test",
            },
            enableColor: true,
            levelColors: map[level.Level]color.Color{
                level.Debug: color.White,
                level.Info:  color.Green,
                level.Warn:  color.Yellow,
                level.Error: color.Red,
                level.Panic: color.Magenta,
            },
            want: color.Green.Colorize("test"),
        },
        {
            name: "Colorize fields",
            args: args{
                level: level.Error,
                msg:   "test",
            },
            enableColor: true,
            levelColors: map[level.Level]color.Color{
                level.Debug: color.White,
                level.Info:  color.Green,
                level.Warn:  color.Yellow,
                level.Error: color.Red,
                level.Panic: color.Magenta,
            },
            prefixFields: []field.Field{
                field.NewTagField("tag"),
                field.NewLevelField(bracket.Angle),
            },
            suffixFields: []field.Field{
                field.NewTagField("tag"),
                field.NewLevelField(bracket.Angle),
            },
            want: color.Red.Colorize("[tag] <ERROR> test [tag] <ERROR>"),
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f, err := NewFormatter(tt.prefixFields, tt.suffixFields, tt.levelColors, tt.enableColor)
            if (err != nil) != tt.wantErr {
                t.Errorf("NewFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if err != nil {
                return
            }

            if got := f.Format(tt.args.level, tt.args.msg); got != tt.want {
                t.Errorf("Format() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_ultraFormatter_Formatf(t *testing.T) {
    type args struct {
        level  level.Level
        format string
        args   []any
    }
    tests := []struct {
        name         string
        prefixFields []field.Field
        suffixFields []field.Field
        args         args
        want         string
        wantErr      bool
    }{
        {
            name: "Default",
            args: args{
                level:  level.Info,
                format: "%v %v",
                args:   []any{"test", "test"},
            },
            prefixFields: []field.Field{
                field.NewTagField("tag"),
                field.NewLevelField(bracket.Angle),
            },
            want: "[tag] <INFO> test test",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f, err := NewFormatter(tt.prefixFields, tt.suffixFields, nil, false)
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
        initPrefixFields        []field.Field
        initSuffixFields        []field.Field
        newPrefixFields         []field.Field
        newSuffixFields         []field.Field
        msg                     string
        logLevel                level.Level
        before                  string
        want                    string
        wantSetPrefixFieldError bool
        wantSetSuffixFieldError bool
    }{
        {
            name: "SetPrefixFields",
            initPrefixFields: []field.Field{
                field.NewTagField("tag"),
                field.NewLevelField(bracket.Angle),
            },
            initSuffixFields: []field.Field{},
            newPrefixFields: []field.Field{
                field.NewTagField("newTag", field.WithBracket(bracket.Round)),
            },
            msg:      "test",
            logLevel: level.Info,
            before:   "[tag] <INFO> test",
            want:     "(newTag) test",
        },
        {
            name:             "SetSuffixFields",
            initPrefixFields: []field.Field{},
            initSuffixFields: []field.Field{
                field.NewTagField("tag"),
                field.NewLevelField(bracket.Angle),
            },
            newPrefixFields: []field.Field{},
            newSuffixFields: []field.Field{
                field.NewTagField("newTag", field.WithBracket(bracket.Round)),
            },
            msg:      "test",
            logLevel: level.Info,
            before:   "test [tag] <INFO>",
            want:     "test (newTag)",
        },
        {
            name:             "SetPrefixFields with existing suffix fields",
            initPrefixFields: []field.Field{},
            initSuffixFields: []field.Field{
                field.NewTagField("tag"),
                field.NewLevelField(bracket.Angle),
            },
            newPrefixFields: []field.Field{
                field.NewTagField("newTag", field.WithBracket(bracket.Round)),
            },
            newSuffixFields: nil,
            msg:             "test",
            logLevel:        level.Info,
            before:          "test [tag] <INFO>",
            want:            "(newTag) test [tag] <INFO>",
        },
        {
            name: "SetSuffixFields with existing prefix fields",
            initPrefixFields: []field.Field{
                field.NewTagField("tag"),
                field.NewLevelField(bracket.Angle),
            },
            initSuffixFields: []field.Field{},
            newPrefixFields:  nil,
            newSuffixFields: []field.Field{
                field.NewTagField("newTag", field.WithBracket(bracket.Round)),
            },
            msg:      "test",
            logLevel: level.Info,
            before:   "[tag] <INFO> test",
            want:     "[tag] <INFO> test (newTag)",
        },
        {
            name: "SetPrefixField with invalid field throws error",
            initPrefixFields: []field.Field{
                field.NewTagField("tag"),
                field.NewLevelField(bracket.Angle),
            },
            initSuffixFields: []field.Field{},
            newPrefixFields: []field.Field{
                invalidField{},
            },
            msg:                     "test",
            logLevel:                level.Info,
            before:                  "[tag] <INFO> test",
            want:                    "[tag] <INFO> test",
            wantSetPrefixFieldError: true,
        },
        {
            name:             "SetSuffixField with invalid field throws error",
            initPrefixFields: []field.Field{},
            initSuffixFields: []field.Field{
                field.NewTagField("tag"),
                field.NewLevelField(bracket.Angle),
            },
            newSuffixFields: []field.Field{
                invalidField{},
            },
            msg:                     "test",
            logLevel:                level.Info,
            before:                  "test [tag] <INFO>",
            wantSetSuffixFieldError: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            f, _ := NewFormatter(tt.initPrefixFields, tt.initSuffixFields, nil, false)

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
