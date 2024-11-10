package field

import (
    "github.com/fmdunlap/go-ultralogger/v2/bracket"
    "testing"
)

func TestTagField_FieldPrinter(t *testing.T) {
    tests := []struct {
        name      string
        tagField  Field
        printArgs PrintArgs
        want      string
        wantErr   bool
    }{
        {
            name:      "Default",
            tagField:  NewTagField("test"),
            printArgs: PrintArgs{},
            want:      "[test]",
        },
        {
            name:      "With Bracket",
            tagField:  NewTagField("test", WithBracket(bracket.Round)),
            printArgs: PrintArgs{},
            want:      "(test)",
        },
        {
            name:      "With Padding",
            tagField:  NewTagField("test", WithPadSettings(TagPadSettings{PadChar: "-", PrefixPadSize: 1, SuffixPadSize: 2})),
            printArgs: PrintArgs{},
            want:      "-[test]--",
        },
        {
            name:      "With Prefix Pad",
            tagField:  NewTagField("test", WithPrefixPadSize(5)),
            printArgs: PrintArgs{},
            want:      "     [test]",
        },
        {
            name:      "With Suffix Pad",
            tagField:  NewTagField("test", WithSuffixPadSize(5)),
            printArgs: PrintArgs{},
            want:      "[test]     ",
        },
        {
            name:      "With Pad Char",
            tagField:  NewTagField("test", WithPrefixPadSize(5), WithPadChar("!")),
            printArgs: PrintArgs{},
            want:      "!!!!![test]",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := tt.tagField.FieldPrinter()
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
