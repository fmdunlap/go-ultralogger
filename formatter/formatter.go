package formatter

import (
    "fmt"
    "github.com/fmdunlap/go-ultralogger/color"
    "github.com/fmdunlap/go-ultralogger/field"
    "github.com/fmdunlap/go-ultralogger/level"
    "maps"
    "strings"
)

type Formatter interface {
    Format(level level.Level, msg string) string
    Formatf(level level.Level, format string, args ...any) string
    SetPrefixFields(fields ...field.Field) error
    SetSuffixFields(fields ...field.Field) error
}

type ColorizedFormatter interface {
    Formatter
    EnableColorization(colorize bool) error
    SetLevelColors(colors map[level.Level]color.Color) error
}

var defaultLevelColors = map[level.Level]color.Color{
    level.Debug: color.Green,
    level.Info:  color.White,
    level.Warn:  color.Yellow,
    level.Error: color.Red,
    level.Panic: color.Magenta,
}

func NewColorizedFormatter(
    prefixFields []field.Field,
    suffixFields []field.Field,
    enableColor bool,
) (ColorizedFormatter, error) {
    levelColors := make(map[level.Level]color.Color)
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
    prefixFields []field.Field

    // suffixFields are the fields that will be printed after the log message.
    suffixFields []field.Field

    // prefixFieldPrintFuncs is an array, in order, of the prefix field printFuncs.
    prefixFieldPrinterFuncs []field.FieldPrinterFunc

    // suffixFieldPrintFuncs is an array, in order, of the suffix field printFuncs.
    suffixFieldPrinterFuncs []field.FieldPrinterFunc

    // colorize indicates whether the log message should be colorized. If true, the log message will be colorized. If
    // false, the log message will not be colorized.
    //
    // TODO: Make colorization configurable via a field color mask or similar.
    colorize bool

    // levelColors is a map of log levels to colors. It is used to colorize the log message based on the log level.
    levelColors map[level.Level]color.Color
}

func (f *ultraFormatter) Format(level level.Level, msg string) string {
    fieldArgs := field.PrintArgs{
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

func (f *ultraFormatter) Formatf(level level.Level, format string, args ...any) string {
    return f.Format(level, fmt.Sprintf(format, args...))
}

func (f *ultraFormatter) SetPrefixFields(fields ...field.Field) error {
    f.prefixFields = fields
    return f.updatePrinterFuncs()
}

func (f *ultraFormatter) SetSuffixFields(fields ...field.Field) error {
    f.suffixFields = fields
    return f.updatePrinterFuncs()
}

func (f *ultraFormatter) updatePrinterFuncs() error {
    f.prefixFieldPrinterFuncs = make([]field.FieldPrinterFunc, len(f.prefixFields))
    for i, fld := range f.prefixFields {
        printerFunc, err := fld.FieldPrinter()
        if err != nil {
            return err
        }
        f.prefixFieldPrinterFuncs[i] = printerFunc
    }

    f.suffixFieldPrinterFuncs = make([]field.FieldPrinterFunc, len(f.suffixFields))
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

func (f *ultraFormatter) SetLevelColors(colors map[level.Level]color.Color) error {
    for _, level := range level.AllLevels() {
        if _, ok := colors[level]; ok {
            f.levelColors[level] = colors[level]
        }
    }
    return nil
}
