package ultralogger

import (
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
            levelField: NewLevelField(BracketAngle),
            printArgs: PrintArgs{
                Level: Info,
            },
            want: "<INFO>",
        },
        {
            name:       "Round Bracket",
            levelField: NewLevelField(BracketRound),
            printArgs: PrintArgs{
                Level: Info,
            },
            want: "(INFO)",
        },
        {
            name:       "Debug",
            levelField: NewLevelField(BracketAngle),
            printArgs: PrintArgs{
                Level: Debug,
            },
            want: "<DEBUG>",
        },
        {
            name:       "Warn",
            levelField: NewLevelField(BracketAngle),
            printArgs: PrintArgs{
                Level: Warn,
            },
            want: "<WARN>",
        },
        {
            name:       "Error",
            levelField: NewLevelField(BracketAngle),
            printArgs: PrintArgs{
                Level: Error,
            },
            want: "<ERROR>",
        },
        {
            name:       "Panic",
            levelField: NewLevelField(BracketAngle),
            printArgs: PrintArgs{
                Level: Panic,
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
