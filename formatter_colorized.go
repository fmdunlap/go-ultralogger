package ultralogger

import (
    "maps"
)

var defaultLevelColors = map[Level]Color{
    Debug: ColorGreen,
    Info:  ColorWhite,
    Warn:  ColorYellow,
    Error: ColorRed,
    Panic: ColorMagenta,
}

type ColorizedFormatter struct {
    BaseFormatter LogLineFormatter
    LevelColors   map[Level]Color
}

func (f *ColorizedFormatter) FormatLogLine(args LogLineArgs, data any) FormatResult {
    res := f.BaseFormatter.FormatLogLine(args, data)
    if res.err != nil {
        return res
    }

    color, ok := f.LevelColors[args.Level]
    if !ok {
        return FormatResult{res.bytes, &MissingLevelColorError{level: args.Level}}
    }

    return FormatResult{color.Colorize(res.bytes), nil}
}

func NewColorizedFormatter(baseFormatter LogLineFormatter, levelColors map[Level]Color) *ColorizedFormatter {
    if levelColors == nil {
        levelColors = make(map[Level]Color)
        maps.Copy(levelColors, defaultLevelColors)
    }

    return &ColorizedFormatter{
        BaseFormatter: baseFormatter,
        LevelColors:   levelColors,
    }
}
