package ultralogger

import (
    "testing"
)

func TestLevelField_FieldPrinter(t *testing.T) {
    tests := []struct {
        name       string
        levelField Field
        llCtx      LogLineContext
        want       string
        wantErr    bool
    }{
        {
            name:       "Default",
            levelField: NewLevelField(BracketAngle),
            llCtx: LogLineContext{
                Level: Info,
            },
            want: "<INFO>",
        },
        {
            name:       "Round Bracket",
            levelField: NewLevelField(BracketRound),
            llCtx: LogLineContext{
                Level: Info,
            },
            want: "(INFO)",
        },
        {
            name:       "Debug",
            levelField: NewLevelField(BracketAngle),
            llCtx: LogLineContext{
                Level: Debug,
            },
            want: "<DEBUG>",
        },
        {
            name:       "Warn",
            levelField: NewLevelField(BracketAngle),
            llCtx: LogLineContext{
                Level: Warn,
            },
            want: "<WARN>",
        },
        {
            name:       "Error",
            levelField: NewLevelField(BracketAngle),
            llCtx: LogLineContext{
                Level: Error,
            },
            want: "<ERROR>",
        },
        {
            name:       "Panic",
            levelField: NewLevelField(BracketAngle),
            llCtx: LogLineContext{
                Level: Panic,
            },
            want: "<PANIC>",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := tt.levelField.FieldFormatter()
            if (err != nil) != tt.wantErr {
                t.Errorf("FieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got(tt.llCtx, OutputFormatText, nil).Data != tt.want {
                t.Errorf("FieldFormatter() got = %v, want %v", got(tt.llCtx, OutputFormatText, nil).Data, tt.want)
            }
        })
    }
}
