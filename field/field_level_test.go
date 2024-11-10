package field

import (
    "github.com/fmdunlap/go-ultralogger/bracket"
    "github.com/fmdunlap/go-ultralogger/level"
    "testing"
)

func TestLevelField_FieldPrinter(t *testing.T) {
    tests := []struct {
        name       string
        levelField Field
        printArgs  PrintArgs
        want       string
        wantErr    bool
    }{
        {
            name:       "Default",
            levelField: NewLevelField(bracket.Angle),
            printArgs: PrintArgs{
                Level: level.Info,
            },
            want: "<INFO>",
        },
        {
            name:       "Round Bracket",
            levelField: NewLevelField(bracket.Round),
            printArgs: PrintArgs{
                Level: level.Info,
            },
            want: "(INFO)",
        },
        {
            name:       "Debug",
            levelField: NewLevelField(bracket.Angle),
            printArgs: PrintArgs{
                Level: level.Debug,
            },
            want: "<DEBUG>",
        },
        {
            name:       "Warn",
            levelField: NewLevelField(bracket.Angle),
            printArgs: PrintArgs{
                Level: level.Warn,
            },
            want: "<WARN>",
        },
        {
            name:       "Error",
            levelField: NewLevelField(bracket.Angle),
            printArgs: PrintArgs{
                Level: level.Error,
            },
            want: "<ERROR>",
        },
        {
            name:       "Panic",
            levelField: NewLevelField(bracket.Angle),
            printArgs: PrintArgs{
                Level: level.Panic,
            },
            want: "<PANIC>",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := tt.levelField.FieldPrinter()
            if (err != nil) != tt.wantErr {
                t.Errorf("FieldPrinter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got(tt.printArgs) != tt.want {
                t.Errorf("FieldPrinter() got = %v, want %v", got(tt.printArgs), tt.want)
            }
        })
    }
}
