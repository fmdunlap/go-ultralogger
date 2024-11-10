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

func (f *FieldLevel) FieldPrinter() (FieldPrinterFunc, error) {
    if f.levelStrings == nil {
        f.levelStrings = make(map[Level]string)

        for _, lvl := range AllLevels() {
            f.levelStrings[lvl] = f.bracket.Wrap(lvl.String())
        }
    }

    return func(args PrintArgs) string {
        return f.levelStrings[args.Level]
    }, nil
}
