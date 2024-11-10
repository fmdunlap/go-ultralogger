package field

import (
    "github.com/fmdunlap/go-ultralogger/v2/bracket"
    "github.com/fmdunlap/go-ultralogger/v2/level"
)

func NewLevelField(bracket bracket.Bracket) *LevelField {
    return &LevelField{
        bracket: bracket,
    }
}

type LevelField struct {
    bracket      bracket.Bracket
    levelStrings map[level.Level]string
}

func (f *LevelField) FieldPrinter() (FieldPrinterFunc, error) {
    if f.levelStrings == nil {
        f.levelStrings = make(map[level.Level]string)

        for _, lvl := range level.AllLevels() {
            f.levelStrings[lvl] = f.bracket.Wrap(lvl.String())
        }
    }

    return func(args PrintArgs) string {
        return f.levelStrings[args.Level]
    }, nil
}
