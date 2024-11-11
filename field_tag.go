package ultralogger

import (
    "strings"
)

type FieldTag struct {
    bracket     Bracket
    padSettings *TagPadSettings
}

type TagPadSettings struct {
    PadChar       string
    PrefixPadSize int
    SuffixPadSize int
}

func NewTagField(opts ...TagFieldOpt) *FieldTag {
    tf := &FieldTag{
        bracket: BracketSquare,
        padSettings: &TagPadSettings{
            PadChar:       " ",
            PrefixPadSize: 0,
            SuffixPadSize: 0,
        },
    }

    for _, opt := range opts {
        opt(tf)
    }

    return tf
}

func (f *FieldTag) FieldFormatter() (FieldFormatter, error) {
    return f.format, nil
}

func (f *FieldTag) format(mCtx LogLineContext, outputFormat OutputFormat, _ any) *FieldResult {
    result := &FieldResult{
        Name: "tag",
    }

    switch outputFormat {
    case OutputFormatJSON:
        result.Data = mCtx.Tag
    case OutputFormatText:
        result.Data = f.tagString(mCtx.Tag)
    }

    return result
}

type TagFieldOpt func(tf *FieldTag)

func WithPadSettings(padSettings TagPadSettings) TagFieldOpt {
    return func(tf *FieldTag) {
        tf.padSettings = &padSettings
    }
}

func WithPadChar(padChar string) TagFieldOpt {
    return func(tf *FieldTag) {
        tf.padSettings.PadChar = padChar
    }
}

func WithPrefixPadSize(prefixPadSize int) TagFieldOpt {
    return func(tf *FieldTag) {
        tf.padSettings.PrefixPadSize = prefixPadSize
    }
}

func WithSuffixPadSize(suffixPadSize int) TagFieldOpt {
    return func(tf *FieldTag) {
        tf.padSettings.SuffixPadSize = suffixPadSize
    }
}

func WithBracket(bracket Bracket) TagFieldOpt {
    return func(tf *FieldTag) {
        tf.bracket = bracket
    }
}

func (f *FieldTag) tagString(tag string) string {
    if tag == "" {
        return ""
    }

    b := strings.Builder{}

    b.WriteString(strings.Repeat(f.padSettings.PadChar, f.padSettings.PrefixPadSize))

    b.WriteString(f.bracket.Open())
    b.WriteString(tag)
    b.WriteString(f.bracket.Close())

    b.WriteString(strings.Repeat(f.padSettings.PadChar, f.padSettings.SuffixPadSize))

    return b.String()
}
