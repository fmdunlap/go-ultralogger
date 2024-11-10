package ultralogger

import (
    "fmt"
    "maps"
    "strings"
)

type Formatter interface {
    Format(level Level, msg string) string
    Formatf(level Level, format string, args ...any) string
    SetPrefixFields(fields ...Field) error
    SetSuffixFields(fields ...Field) error
}

type ColorizedFormatter interface {
    Formatter
    EnableColorization(colorize bool) error
    SetLevelColors(colors map[Level]Color) error
}

var defaultLevelColors = map[Level]Color{
    Debug: ColorGreen,
    Info:  ColorWhite,
    Warn:  ColorYellow,
    Error: ColorRed,
    Panic: ColorMagenta,
}

func NewColorizedFormatter(
    prefixFields []Field,
    suffixFields []Field,
    enableColor bool,
) (ColorizedFormatter, error) {
    levelColors := make(map[Level]Color)
    maps.Copy(levelColors, defaultLevelColors)

    uf := &ultraFormatter{
        prefixFields: prefixFields,
        suffixFields: suffixFields,
        levelColors:  levelColors,
        colorize:     enableColor,
    }

    err := uf.updatePrinterFuncs()
    if err != nil {
        return nil, err
    }

    return uf, nil
}

type ultraFormatter struct {
    // prefixFields are the fields that will be printed before the log message.
    prefixFields []Field

    // suffixFields are the fields that will be printed after the log message.
    suffixFields []Field

    // prefixFieldPrintFuncs is an array, in order, of the prefix field printFuncs.
    prefixFieldPrinterFuncs []FieldPrinterFunc

    // suffixFieldPrintFuncs is an array, in order, of the suffix field printFuncs.
    suffixFieldPrinterFuncs []FieldPrinterFunc

    // colorize indicates whether the log message should be colorized. If true, the log message will be colorized. If
    // false, the log message will not be colorized.
    //
    // TODO: Make colorization configurable via a field color mask or similar.
    colorize bool

    // levelColors is a map of log levels to colors. It is used to colorize the log message based on the log level.
    levelColors map[Level]Color
}

func (f *ultraFormatter) Format(level Level, msg string) string {
    fieldArgs := PrintArgs{
        Level: level,
    }

    b := strings.Builder{}

    for _, printerFunc := range f.prefixFieldPrinterFuncs {
        b.WriteString(printerFunc(fieldArgs))
        b.WriteRune(' ')
    }

    b.WriteString(msg)

    for _, printerFunc := range f.suffixFieldPrinterFuncs {
        b.WriteRune(' ')
        b.WriteString(printerFunc(fieldArgs))
    }

    if f.colorize {
        return f.levelColors[level].Colorize(b.String())
    }

    return b.String()
}

func (f *ultraFormatter) Formatf(level Level, format string, args ...any) string {
    return f.Format(level, fmt.Sprintf(format, args...))
}

func (f *ultraFormatter) SetPrefixFields(fields ...Field) error {
    f.prefixFields = fields
    return f.updatePrinterFuncs()
}

func (f *ultraFormatter) SetSuffixFields(fields ...Field) error {
    f.suffixFields = fields
    return f.updatePrinterFuncs()
}

func (f *ultraFormatter) updatePrinterFuncs() error {
    f.prefixFieldPrinterFuncs = make([]FieldPrinterFunc, len(f.prefixFields))
    for i, fld := range f.prefixFields {
        printerFunc, err := fld.FieldPrinter()
        if err != nil {
            return err
        }
        f.prefixFieldPrinterFuncs[i] = printerFunc
    }

    f.suffixFieldPrinterFuncs = make([]FieldPrinterFunc, len(f.suffixFields))
    for i, fld := range f.suffixFields {
        printerFunc, err := fld.FieldPrinter()
        if err != nil {
            return err
        }
        f.suffixFieldPrinterFuncs[i] = printerFunc
    }

    return nil
}

func (f *ultraFormatter) EnableColorization(colorize bool) error {
    f.colorize = colorize
    return nil
}

func (f *ultraFormatter) SetLevelColors(colors map[Level]Color) error {
    for _, level := range AllLevels() {
        if _, ok := colors[level]; ok {
            f.levelColors[level] = colors[level]
        }
    }
    return nil
}
