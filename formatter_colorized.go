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

func (f *ColorizedFormatter) FormatLogLine(mCtx LogLineContext, data any) ([]byte, error) {
    logLine, err := f.BaseFormatter.FormatLogLine(mCtx, data)
    if err != nil {
        return nil, err
    }

    color, ok := f.LevelColors[mCtx.Level]
    if !ok {
        return logLine, &MissingLevelColorError{level: mCtx.Level}
    }

    return color.Colorize(logLine), nil
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
