package field

import (
    "github.com/fmdunlap/go-ultralogger/v2/bracket"
    "strings"
)

type TagField struct {
    tag            string
    bracket        bracket.Bracket
    padSettings    *TagPadSettings
    precomputedTag *string
}

type TagPadSettings struct {
    PadChar       string
    PrefixPadSize int
    SuffixPadSize int
}

func NewTagField(tag string, opts ...TagFieldOpt) *TagField {
    tf := &TagField{
        tag:     tag,
        bracket: bracket.Square,
        padSettings: &TagPadSettings{
            PadChar:       " ",
            PrefixPadSize: 0,
            SuffixPadSize: 0,
        },
    }

    for _, opt := range opts {
        opt(tf)
    }

    tf.updatePrecomputedTag()

    return tf
}

func (f *TagField) FieldPrinter() (FieldPrinterFunc, error) {
    if f.precomputedTag == nil {
        f.updatePrecomputedTag()
    }

    return func(info PrintArgs) string {
        return *f.precomputedTag
    }, nil
}

type TagFieldOpt func(tf *TagField)

func WithPadSettings(padSettings TagPadSettings) TagFieldOpt {
    return func(tf *TagField) {
        tf.padSettings = &padSettings
    }
}

func WithPadChar(padChar string) TagFieldOpt {
    return func(tf *TagField) {
        tf.padSettings.PadChar = padChar
    }
}

func WithPrefixPadSize(prefixPadSize int) TagFieldOpt {
    return func(tf *TagField) {
        tf.padSettings.PrefixPadSize = prefixPadSize
    }
}

func WithSuffixPadSize(suffixPadSize int) TagFieldOpt {
    return func(tf *TagField) {
        tf.padSettings.SuffixPadSize = suffixPadSize
    }
}

func WithBracket(bracket bracket.Bracket) TagFieldOpt {
    return func(tf *TagField) {
        tf.bracket = bracket
    }
}

func (f *TagField) updatePrecomputedTag() {
    tagStr := f.tagString()
    f.precomputedTag = &tagStr
}

func (f *TagField) tagString() string {
    if f.tag == "" {
        return ""
    }

    b := strings.Builder{}

    b.WriteString(strings.Repeat(f.padSettings.PadChar, f.padSettings.PrefixPadSize))

    b.WriteString(f.bracket.Open())
    b.WriteString(f.tag)
    b.WriteString(f.bracket.Close())

    b.WriteString(strings.Repeat(f.padSettings.PadChar, f.padSettings.SuffixPadSize))

    return b.String()
}
