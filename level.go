package ultralogger

import (
    "strings"
)

// Level is a type representing the level of a log message.
//
// It can be one of the following:
//   - DebugLevel
//   - InfoLevel
//   - WarnLevel
//   - ErrorLevel
//   - PanicLevel
//
// Levels determine the priority of a log message, and can be hidden if a logger's minimum level is set to a higher
// level than the message's level.
//
// For example, if a logger's minimum level is set to WarnLevel, then a message with a level of InfoLevel will not be
// written to the output.
type Level int

var defaultLevelColors = map[Level]Color{
    DebugLevel: ColorGreen,
    InfoLevel:  ColorWhite,
    WarnLevel:  ColorYellow,
    ErrorLevel: ColorRed,
    PanicLevel: ColorMagenta,
}

const (
    DebugLevel Level = iota
    InfoLevel
    WarnLevel
    ErrorLevel
    PanicLevel
)

// AllLevels returns a slice of all available levels.
func AllLevels() []Level {
    return []Level{
        DebugLevel,
        InfoLevel,
        WarnLevel,
        ErrorLevel,
        PanicLevel,
    }
}

func (l Level) String() string {
    switch l {
    case DebugLevel:
        return "DEBUG"
    case InfoLevel:
        return "INFO"
    case WarnLevel:
        return "WARN"
    case ErrorLevel:
        return "ERROR"
    case PanicLevel:
        return "PANIC"
    default:
        return "UNKNOWN"
    }
}

// ParseLevel parses a string into a Level. Returns an error if the string is not a valid Level.
func ParseLevel(levelStr string) (Level, error) {
    switch strings.ToLower(levelStr) {
    case "debug":
        return DebugLevel, nil
    case "info":
        return InfoLevel, nil
    case "warn":
        return WarnLevel, nil
    case "error":
        return ErrorLevel, nil
    case "panic":
        return PanicLevel, nil
    default:
        return 0, &LevelParsingError{level: levelStr}
    }
}
