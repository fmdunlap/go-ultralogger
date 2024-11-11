package ultralogger

func NewLevelField(bracket Bracket) *FieldLevel {
    return &FieldLevel{
        bracket: bracket,
    }
}

type FieldLevel struct {
    bracket      Bracket
    levelStrings map[Level]string
}

func (f *FieldLevel) FieldFormatter() (FieldFormatter, error) {
    if f.levelStrings == nil {
        f.levelStrings = make(map[Level]string)

        for _, lvl := range AllLevels() {
            f.levelStrings[lvl] = f.bracket.Wrap(lvl.String())
        }
    }

    return f.format, nil
}

func (f *FieldLevel) format(mCtx LogLineContext, outputFormat OutputFormat, _ any) *FieldResult {
    result := &FieldResult{
        Name: "Level",
    }

    switch outputFormat {
    case OutputFormatJSON:
        result.Data = f.levelStrings[mCtx.Level]
    case OutputFormatText:
        result.Data = f.levelStrings[mCtx.Level]
    }

    return result
}
