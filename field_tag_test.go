package ultralogger

import (
    "testing"
)

func TestTagField_FieldPrinter(t *testing.T) {
    tests := []struct {
        name     string
        tagField Field
        llCtx    LogLineContext
        want     string
        wantErr  bool
    }{
        {
            name:     "Default",
            tagField: NewTagField(),
            llCtx: LogLineContext{
                Level: Info,
                Tag:   "test",
            },
            want: "[test]",
        },
        {
            name:     "With Bracket",
            tagField: NewTagField(WithBracket(BracketRound)),
            llCtx: LogLineContext{
                Level: Info,
                Tag:   "test",
            },
            want: "(test)",
        },
        {
            name:     "With Padding",
            tagField: NewTagField(WithPadSettings(TagPadSettings{PadChar: "-", PrefixPadSize: 1, SuffixPadSize: 2})),
            llCtx: LogLineContext{
                Level: Info,
                Tag:   "test",
            },
            want: "-[test]--",
        },
        {
            name:     "With Prefix Pad",
            tagField: NewTagField(WithPrefixPadSize(5)),
            llCtx: LogLineContext{
                Level: Info,
                Tag:   "test",
            },
            want: "     [test]",
        },
        {
            name:     "With Suffix Pad",
            tagField: NewTagField(WithSuffixPadSize(5)),
            llCtx: LogLineContext{
                Level: Info,
                Tag:   "test",
            },
            want: "[test]     ",
        },
        {
            name:     "With Pad Char",
            tagField: NewTagField(WithPrefixPadSize(5), WithPadChar("!")),
            llCtx: LogLineContext{
                Level: Info,
                Tag:   "test",
            },
            want: "!!!!![test]",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := tt.tagField.FieldFormatter()
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
