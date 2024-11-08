package ultralogger

import (
    "fmt"
    "io"
    "maps"
    "os"
    "strings"
    "time"
)

const (
    defaultTagSpaceLen       = 10
    defaultTagBracketType    = BracketTypeSquare
    defaultLevelbracketType  = BracketTypeAngle
    defaultDateFormat        = "2006-01-02"
    defaultTimeFormat        = "15:04:05"
    defaultDateTimeSeparator = " "
)

// UltraLogger is a versatile, flexible, and efficient logging library for Go that supports multiple log levels, custom
// formatting, and output destinations. It is designed to be easy to use while providing advanced features such as
// colorized output, customizable log levels, and detailed error handling.
type UltraLogger struct {
    // minLevel is the minimum log level that will be logged. All log messages with a level lower than this will be
    // ignored.
    minLevel Level

    // levelBracketType is the type of bracket to use for the level in the log message. It can be either BracketTypeAngle
    // or BracketTypeSquare.
    levelBracketType BracketType

    // tag is the tag to be prepended to the log message. It can be used to identify the source of the log message.
    tag string

    // tagBracketType is the type of bracket to use for the tag in the log message. See BracketType for more information
    tagBracketType BracketType

    // padTag indicates whether the tag should be padded with spaces to the specified tagPadSize. If true, the tag will
    // be padded with spaces to the specified tagPadSize. If false, the tag will not be padded.
    padTag bool

    // tagPadSize is the size of the padding to be used for the tag. If padTag is true, the tag will be padded with
    // spaces to this size.
    tagPadSize int

    // showDate indicates whether the date should be included in the log message. If true, the date will be included in
    // the log message. If false, the date will not be included.
    showDate bool

    // dateFormat is the format string for the date to be included in the log message. It can be any valid format string
    // for time.Format.
    dateFormat string

    // showTime indicates whether the time should be included in the log message. If true, the time will be included in
    // the log message. If false, the time will not be included.
    showTime bool

    // timeFormat is the format string for the time to be included in the log message. It can be any valid format string
    // for time.Format.
    timeFormat string

    // dateTimeSeparator is the separator to be used between the date and time in the log message. It can be any string
    // that is not a valid format string for time.Format.
    dateTimeSeparator string

    // colorize indicates whether the log message should be colorized. If true, the log message will be colorized. If
    // false, the log message will not be colorized.
    colorize bool

    // levelColors is a map of log levels to colors. It is used to colorize the log message based on the log level.
    levelColors map[Level]Color

    // silent indicates whether the log message should be printed to the output. If true, the log message will not be
    // printed to the output. If false, the log message will be printed to the output.
    silent bool

    // fallback indicates whether to fallback to the original writer if an error occurs while writing to the output.
    // If true, the original writer will be used if an error occurs while writing to the output. If false, the error
    // will be panicked.
    fallback bool
    // panicOnPanicLevel indicates whether to panic if a panic occurs while logging a message with a panic level.
    // If true, a panic will be panicked if a panic occurs while logging a message with a panic level. If false, the
    // panic will be ignored.
    panicOnPanicLevel bool

    // computedDateTimeFormat is the computed date time format string based on the showDate, showTime, and
    // dateTimeSeparator fields. It is updated whenever settings that affect the format of the log string change via the
    // updateFormatStrings method.
    computedDateTimeFormat string

    // computedFmtString is the computed format string based on the computedDateTimeFormat, levelBracketType, and
    // tag fields. It is updated whenever settings that affect the format of the log string change via the
    // updateFormatStrings method.
    computedFmtString string

    // clock is the clock used to get the current time for the log messages.
    clock clock

    // writer is the writer to which the log messages will be written.
    writer io.Writer
}

// NewUltraLogger creates a new UltraLogger with default settings and the given writer.
//
// Default settings include:
//   - MinLevel: InfoLevel
//   - LevelBracketType: BracketTypeAngle
//   - Tag: ""
//   - TagBracketType: BracketTypeAngle
//   - PadTag: true
//   - TagPadSize: 10
//   - ShowDate: true
//   - DateFormat: "2006-01-02"
//   - ShowTime: true
//   - TimeFormat: "15:04:05"
//   - DateTimeSeparator: " "
//   - Colorize: false
//   - LevelColors: defaultLevelColors
//   - Silent: false
//   - Fallback: true
//   - PanicOnPanicLevel: false
//   - Clock: realClock{}
//   - Writer: writer
func NewUltraLogger(writer io.Writer) *UltraLogger {
    levelColors := make(map[Level]Color)
    maps.Copy(levelColors, defaultLevelColors)

    l := &UltraLogger{
        minLevel:         InfoLevel,
        levelBracketType: defaultLevelbracketType,

        tag:            "",
        tagBracketType: defaultTagBracketType,
        padTag:         true,
        tagPadSize:     defaultTagSpaceLen,

        showDate:          true,
        dateFormat:        defaultDateFormat,
        showTime:          true,
        timeFormat:        defaultTimeFormat,
        dateTimeSeparator: defaultDateTimeSeparator,

        colorize:    false,
        levelColors: levelColors,

        silent: false,

        fallback:          true,
        panicOnPanicLevel: false,

        clock: &realClock{},

        writer: writer,
    }

    l.updateFormatStrings()

    return l
}

// == Logging == //

// Log logs a message with the given level and message.
func (l *UltraLogger) Log(level Level, msg string) {
    if !l.silent && level >= l.minLevel {
        if _, err := l.writer.Write([]byte(l.Slogln(level, msg))); err != nil {
            l.handleLogWriterError(level, msg, err)
        }
    }
}

// Logf logs a formatted message with the given level and format string.
func (l *UltraLogger) Logf(level Level, format string, args ...any) {
    l.Log(level, fmt.Sprintf(format, args...))
}

// Debug logs a message with the DebugLevel level and message.
func (l *UltraLogger) Debug(msg string) {
    l.Log(DebugLevel, msg)
}

// Debugf logs a formatted message with the DebugLevel level and format string.
func (l *UltraLogger) Debugf(format string, args ...any) {
    l.Logf(DebugLevel, format, args...)
}

// Info logs a message with the InfoLevel level and message.
func (l *UltraLogger) Info(msg string) {
    l.Log(InfoLevel, msg)
}

// Infof logs a formatted message with the InfoLevel level and format string.
func (l *UltraLogger) Infof(format string, args ...any) {
    l.Logf(InfoLevel, format, args...)
}

// Warn logs a message with the WarnLevel level and message.
func (l *UltraLogger) Warn(msg string) {
    l.Log(WarnLevel, msg)
}

// Warnf logs a formatted message with the WarnLevel level and format string.
func (l *UltraLogger) Warnf(format string, args ...any) {
    l.Logf(WarnLevel, format, args...)
}

// Error logs a message with the ErrorLevel level and message.
func (l *UltraLogger) Error(msg string) {
    l.Log(ErrorLevel, msg)
}

// Errorf logs a formatted message with the ErrorLevel level and format string.
func (l *UltraLogger) Errorf(format string, args ...any) {
    l.Logf(ErrorLevel, format, args...)
}

// Panic logs a message with the PanicLevel level and message. If panicOnPanicLevel is true, it panics.
func (l *UltraLogger) Panic(msg string) {
    l.Log(PanicLevel, msg)

    if l.panicOnPanicLevel {
        panic(msg)
    }
}

// Panicf logs a formatted message with the PanicLevel level and format string. If panicOnPanicLevel is true, it panics.
func (l *UltraLogger) Panicf(format string, args ...any) {
    l.Logf(PanicLevel, format, args...)

    if l.panicOnPanicLevel {
        panic(fmt.Sprintf(format, args...))
    }
}

// == Settings == //

// SetMinLevel sets the minimum log level that will be logged. All log messages with a level lower than this will be
// ignored. Returns the logger.
func (l *UltraLogger) SetMinLevel(level Level) Logger {
    l.minLevel = level
    return l
}

// GetMinLevel returns the minimum log level that will be logged.
func (l *UltraLogger) GetMinLevel() Level {
    return l.minLevel
}

// SetLevelBracketType sets the type of bracket to use for the level in the log message. It can be any BracketType.
// Returns the logger and an error if the bracket type is invalid.
// BracketTypeAngle - "<level>"
// BracketTypeSquare - "[level]"
// BracketTypeRound - "(level)"
// BracketTypeCurly - "{level}"
// BracketTypeNone - "level"
func (l *UltraLogger) SetLevelBracketType(bracketType BracketType) (Logger, error) {
    if !validBracketType(bracketType) {
        return l, &InvalidBracketTypeError{bracketType: bracketType}
    }
    l.levelBracketType = bracketType
    l.updateFormatStrings()
    return l, nil
}

// GetLevelBracketType returns the type of bracket to use for the level in the log message.
func (l *UltraLogger) GetLevelBracketType() BracketType {
    return l.levelBracketType
}

// SetTag sets the tag to be prepended to the log message. It can be used to identify the source of the log message.
// Returns the logger.
func (l *UltraLogger) SetTag(tag string) Logger {
    l.tag = tag
    l.computedFmtString = l.computeFormatStr()
    return l
}

// GetTag returns the tag to be prepended to the log message.
func (l *UltraLogger) GetTag() string {
    return l.tag
}

// SetTagPaddingEnabled enables or disables the padding of the tag with spaces. If true, the tag will be padded with
// spaces to the specified tagPadSize. If false, the tag will not be padded. Returns the logger.
//
// Example:
//   logger.SetTagPaddingEnabled(true).SetTagPadSize(10)
//   logger.SetTagPaddingEnabled(false).SetTagPadSize(10) // This will have no effect
//   logger.SetTagPaddingEnabled(true).SetTagPadSize(0)   // This will have no effect
//
// Note: If the tag + bracket characters exceed the tagPadSize, the tag will push the message one space past the
// tag and bracket characters.
func (l *UltraLogger) SetTagPaddingEnabled(pad bool) Logger {
    l.padTag = pad
    l.updateFormatStrings()
    return l
}

// GetTagPaddingEnabled returns whether the tag is padded with spaces.
func (l *UltraLogger) GetTagPaddingEnabled() bool {
    return l.padTag
}

// SetTagPadSize sets the size of the padding to be used for the tag. If padTag is true, the tag will be padded with
// spaces to this size. Returns the logger.
//
// Example:
//   logger := ultralogger.NewUltraLogger(os.Stdout)
//   logger.SetTag("TAG")
//
//   logger.Info("Hello, World!")   // "2024-11-07 15:04:05 [TAG] <INFO> Hello, World!"
//   logger.Info("Go UltraLogger!") // "2024-11-07 15:04:05 [TAG] <INFO> Go UltraLogger!"
//
//   logger.SetTagPaddingEnabled(true)
//   logger.SetTagPadSize(10)
//
//   logger.Info("Hello, World!")   // "2024-11-07 15:04:05 [TAG]     <INFO> Hello, World!"
//   logger.Info("Go UltraLogger!") // "2024-11-07 15:04:05 [TAG]     <INFO> Go UltraLogger!"
//
//
//   logger.SetTagPadSize(0)
//   logger.Info("Hello, World!") // "2024-11-07 15:04:05 [TAG] <INFO> Hello, World!"
//
// Note: If the tag + bracket characters exceed the tagPadSize, the tag will push the message one space past the
// tag and bracket characters.
//
//   logger.SetTag("IMAVERYLONGTAG")
//   logger.Info("Hello, World!") // "2024-11-07 15:04:05 [IMAVERYLONGTAG] <INFO> Hello, World!"
func (l *UltraLogger) SetTagPadSize(size int) Logger {
    l.tagPadSize = size
    l.updateFormatStrings()
    return l
}

// GetTagPadSize returns the size of the padding to be used for the tag.
func (l *UltraLogger) GetTagPadSize() int {
    return l.tagPadSize
}

// SetTagBracketType sets the type of bracket to use for the tag in the log message. It can be any BracketType.
// Returns the logger and an error if the bracket type is invalid.
//
// BracketTypeAngle - "<tag>"
// BracketTypeSquare - "[tag]"
// BracketTypeRound - "(tag)"
// BracketTypeCurly - "{tag}"
// BracketTypeNone - "tag"
func (l *UltraLogger) SetTagBracketType(bracketType BracketType) (Logger, error) {
    if !validBracketType(bracketType) {
        return l, &InvalidBracketTypeError{bracketType: bracketType}
    }
    l.tagBracketType = bracketType
    l.updateFormatStrings()
    return l, nil
}

// GetTagBracketType returns the type of bracket to use for the tag in the log message.
func (l *UltraLogger) GetTagBracketType() BracketType {
    return l.tagBracketType
}

// ShowTime enables or disables the time to be included in the log message. If true, the time will be included in the
// log message. If false, the time will not be included. Returns the logger.
func (l *UltraLogger) ShowTime(show bool) Logger {
    l.showTime = show
    l.updateFormatStrings()
    return l
}

// GetShowTime returns whether the time is included in the log message.
func (l *UltraLogger) GetShowTime() bool {
    return l.showTime
}

// SetTimeFormat sets the format string for the time to be included in the log message. It can be any valid format
// string for time.Format. Returns the logger and an error if the format is invalid.
//
// Example:
//   logger := ultralogger.NewUltraLogger(os.Stdout)
//
//   logger.Info("Hello, World!") // "2024-11-07 15:04:05 <INFO> Hello, World!"
//
//   logger.SetTimeFormat("15|04|05")
//   logger.Info("Hello, World!") // "2024-11-07 15|04|05 <INFO> Hello, World!"
func (l *UltraLogger) SetTimeFormat(format string) (Logger, error) {
    if !validTimeFormat(format) {
        return l, &InvalidTimeFormatError{format: format}
    }

    l.timeFormat = format
    l.updateFormatStrings()
    return l, nil
}

// GetTimeFormat returns the format string for the time to be included in the log message.
func (l *UltraLogger) GetTimeFormat() string {
    return l.timeFormat
}

// ShowDate enables or disables the date to be included in the log message. If true, the date will be included in the
// log message. If false, the date will not be included. Returns the logger.
func (l *UltraLogger) ShowDate(show bool) Logger {
    l.showDate = show
    l.updateFormatStrings()
    return l
}

// GetShowDate returns whether the date is included in the log message.
func (l *UltraLogger) GetShowDate() bool {
    return l.showDate
}

// SetDateFormat sets the format string for the date to be included in the log message. It can be any valid format
// string for time.Format. Returns the logger and an error if the format is invalid.
//
// Example:
//   logger := ultralogger.NewUltraLogger(os.Stdout)
//
//   logger.Info("Hello, World!") // "2024-11-07 15:04:05 <INFO> Hello, World!"
//
//   logger.SetDateFormat("15/04/05")
//   logger.Info("Hello, World!") // "2024/11/07 15:04:05 <INFO> Hello, World!"
func (l *UltraLogger) SetDateFormat(format string) (Logger, error) {
    if !validDateFormat(format) {
        return l, &InvalidDateFormatError{format: format}
    }

    l.dateFormat = format
    l.updateFormatStrings()
    return l, nil
}

// GetDateFormat returns the format string for the date to be included in the log message.
func (l *UltraLogger) GetDateFormat() string {
    return l.dateFormat
}

// SetDateTimeSeparator sets the separator to be used between the date and time in the log message. It can be any
// string that is not a valid format string for time.Format. Returns the logger.
//
// Example:
//   logger := ultralogger.NewUltraLogger(os.Stdout)
//
//   logger.Info("Hello, World!") // "2024-11-07 15:04:05 <INFO> Hello, World!"
//
//   logger.SetDateTimeSeparator("@")
//   logger.Info("Hello, World!") // "2024/11/07@15:04:05 <INFO> Hello, World!"
func (l *UltraLogger) SetDateTimeSeparator(separator string) Logger {
    l.dateTimeSeparator = separator
    l.updateFormatStrings()
    return l
}

// GetDateTimeSeparator returns the separator to be used between the date and time in the log message.
func (l *UltraLogger) GetDateTimeSeparator() string {
    return l.dateTimeSeparator
}

// SetColorize sets whether the log message should be colorized. If true, the log message will be colorized. If false,
// the log message will not be colorized. Returns the logger.
func (l *UltraLogger) SetColorize(colorize bool) Logger {
    l.colorize = colorize
    return l
}

// GetColorize returns whether the log message should be colorized.
func (l *UltraLogger) GetColorize() bool {
    return l.colorize
}

// SetLevelColor sets the color for a specific log level. It can be any valid color. Returns the logger and an error if
// the color is invalid.
//
// Example:
//   logger := ultralogger.NewUltraLogger(os.Stdout)
//   logger.SetColorize(true)
//
//   logger.Warn("Hello, World!") // "\033[33m<WARN> Hello, World!\033[0m" (Note: this would be yellow in a terminal)
//
//   logger.SetLevelColor(ultralogger.WarnLevel, ultralogger.ColorGreen)
//   logger.Warn("Hello, World!") // "\033[32m<WARN> Hello, World!\033[0m" (Note: this would be green in a terminal)
func (l *UltraLogger) SetLevelColor(level Level, color Color) (Logger, error) {
    if !validColor(color) {
        return l, &InvalidColorError{color: color}
    }

    if _, ok := l.levelColors[level]; !ok {
        return l, &InvalidLevelError{level: int(level)}
    }

    l.levelColors[level] = color
    return l, nil
}

// GetLevelColor returns the color for a specific log level.
func (l *UltraLogger) GetLevelColor(level Level) Color {
    return l.levelColors[level]
}

// SetLevelColors sets the color for log levels specified in the map. Any valid log level can be any valid color. If
// a level is not specified in the map, it will remain unchanged. Returns the logger and an error if the color is
// invalid.
//
// Example:
//   logger := ultralogger.NewUltraLogger(os.Stdout)
//   logger.SetColorize(true)
//
//   logger.Warn("Hello, World!") // "\033[33m<WARN> Hello, World!\033[0m" (Note: this would be yellow in a terminal)
//
//   logger.SetLevelColors(map[ultralogger.Level]ultralogger.Color{
//       ultralogger.WarnLevel: ultralogger.ColorGreen,
//       ultralogger.ErrorLevel: ultralogger.ColorRed,
//   })
//   logger.Warn("Hello, World!") // "\033[32m<WARN> Hello, World!\033[0m" (Note: this would be green in a terminal)
//   logger.Error("Hello, World!") // "\033[31m<ERROR> Hello, World!\033[0m" (Note: this would be red in a terminal)
func (l *UltraLogger) SetLevelColors(colors map[Level]Color) (Logger, error) {
    for level, color := range colors {
        if !validColor(color) {
            return l, &InvalidColorError{color: color}
        }

        _, ok := l.levelColors[level]
        if !ok {
            return l, &InvalidLevelError{level: int(level)}
        }

        l.levelColors[level] = color
    }
    return l, nil
}

// GetLevelColors returns a map of log levels to colors.
func (l *UltraLogger) GetLevelColors() map[Level]Color {
    return l.levelColors
}

// SetSilent sets whether the log message should be printed to the output. If true, the log message will not be
// printed to the output. If false, the log message will be printed to the output. Returns the logger.
//
// Example:
//   logger := ultralogger.NewUltraLogger(os.Stdout)
//
//   logger.Info("Hello, World!") // "2024-11-07 15:04:05 <INFO> Hello, World!"
//   logger.SetSilent(true)
//   logger.Info("Hello, World!") // Nothing will be printed to the output
func (l *UltraLogger) SetSilent(silence bool) Logger {
    l.silent = silence
    return l
}

// GetSilent returns whether the log message should be printed to the output.
func (l *UltraLogger) GetSilent() bool {
    return l.silent
}

// SetPanicOnPanicLevel sets whether to panic if a panic occurs while logging a message with a panic level. If true,
// a panic will be panicked if a panic occurs while logging a message with a panic level. If false, the panic will be
// ignored. Returns the logger.
func (l *UltraLogger) SetPanicOnPanicLevel(shouldPanic bool) Logger {
    l.panicOnPanicLevel = shouldPanic
    return l
}

// GetPanicOnPanicLevel returns whether to panic if a panic occurs while logging a message with a panic level.
func (l *UltraLogger) GetPanicOnPanicLevel() bool {
    return l.panicOnPanicLevel
}

// Writer returns the writer to which the log messages will be written.
func (l *UltraLogger) Writer() io.Writer {
    return l.writer
}

// == Strings and Formatting == //

// Slog returns the string representation of a log message with the given level and message.
func (l *UltraLogger) Slog(level Level, msg string) string {
    return l.format(level, msg)
}

// Slogf returns the string representation of a formatted log message with the given level and format string.
func (l *UltraLogger) Slogf(level Level, format string, args ...any) string {
    return l.format(level, fmt.Sprintf(format, args...))
}

// Slogln returns the string representation of a log message with the given level and message, followed by a newline.
func (l *UltraLogger) Slogln(level Level, msg string) string {
    return l.format(level, msg) + "\n"
}

// updateFormatStrings updates the computed format strings based on the current settings.
func (l *UltraLogger) updateFormatStrings() {
    l.computedFmtString = l.computeFormatStr()
    l.computedDateTimeFormat = l.computeDatetimeFormat()
}

// format returns the string representation of a log message with the given level and message.
// This is the function that ultimately does the formatting of the log message.
func (l *UltraLogger) format(level Level, msg string) string {
    levelString := l.levelBracketType.wrap(level.String())
    if l.colorize {
        levelString = colorize(l.levelColor(level), levelString)
        msg = colorize(l.levelColor(level), msg)
    }

    if l.showDate || l.showTime {
        formattedTime := l.clock.Now().Format(l.computedDateTimeFormat)
        return fmt.Sprintf(l.computedFmtString, formattedTime, levelString, msg)
    }

    return fmt.Sprintf(l.computedFmtString, levelString, msg)
}

// computeFormatStr computes the format string for the logger based on the
// Format: [\[tag\]] [date@time|date|time] <level> message
func (l *UltraLogger) computeFormatStr() string {
    dateTimePart := ""
    if l.showDate || l.showTime {
        dateTimePart = "%s "
    }

    return dateTimePart + l.computeTagPrefix() + "%s %s"
}

// computeTagPrefix creates the precomputed format string for the tag.
func (l *UltraLogger) computeTagPrefix() string {
    if l.tag == "" {
        return ""
    }

    tagPart := l.tagBracketType.wrap(l.tag) + " "
    if l.padTag {
        extraSpace := l.tagPadSize - len(tagPart)
        if extraSpace > 0 {
            tagPart += strings.Repeat(" ", extraSpace)
        }
    }

    return tagPart
}

// computeDatetimeFormat creates the precomputed format string for date and time.
func (l *UltraLogger) computeDatetimeFormat() string {
    if l.showDate && l.showTime {
        return l.dateFormat + l.dateTimeSeparator + l.timeFormat
    } else if l.showTime {
        return l.timeFormat
    } else if l.showDate {
        return l.dateFormat
    }

    return ""
}

// levelColor returns the color for a specific log level.
func (l *UltraLogger) levelColor(level Level) Color {
    return l.levelColors[level]
}

// == Error Handling == //

// handleLogWriterError handles errors that occur while writing to the output. On failure, the log will fallback to
// writing to os.Stdout.
func (l *UltraLogger) handleLogWriterError(level Level, msg string, err error) {
    if !l.fallback || l.writer == os.Stdout {
        panic(err)
    }

    l.writer = os.Stdout
    l.Logf(
        ErrorLevel,
        "error writing to original log writer, falling back to stdout: %v",
        err,
    )
    l.Log(level, msg)
}

// validDateFormat checks if the provided format string is a valid date format.
//
// Note: I'm convinced there's a better way to do this... but I have no idea what it is. Send a PR!
func validDateFormat(format string) bool {
    testDate := time.Date(2006, time.January, 2, 0, 4, 5, 0, time.FixedZone("MST", -7*60*60))
    parsedDate, err := time.Parse(format, testDate.Format(format))
    return err == nil && parsedDate.Year() == testDate.Year() && parsedDate.Month() == testDate.Month() && parsedDate.Day() == testDate.Day()
}

// validTimeFormat checks if the provided format string is a valid time format.
//
// Note: I'm convinced there's a better way to do this... but I have no idea what it is. Send a PR!
func validTimeFormat(format string) bool {
    testTime := time.Date(2006, time.January, 2, 15, 4, 5, 0, time.FixedZone("MST", -7*60*60))
    parsedTime, err := time.Parse(format, testTime.Format(format))
    return err == nil && parsedTime.Hour() == testTime.Hour() && parsedTime.Minute() == testTime.Minute() && parsedTime.Second() == testTime.Second()
}
