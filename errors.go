package ultralogger

import (
    "errors"
    "fmt"
)

// == Level Errors == //

type LevelParsingError struct {
    level string
}

func (e *LevelParsingError) Error() string {
    return fmt.Sprintf("invalid level: %s", e.level)
}

type InvalidLevelError struct {
    level int
}

func (e *InvalidLevelError) Error() string {
    return fmt.Sprintf("invalid level: %d", e.level)
}

type MissingLevelsError struct {
    levels []Level
}

func (m MissingLevelsError) Error() string {
    return fmt.Sprintf("missing levels: %v", m.levels)
}

// == Bracket Errors == //

type InvalidBracketTypeError struct {
    bracketType BracketType
}

func (e *InvalidBracketTypeError) Error() string {
    return fmt.Sprintf("invalid bracket type: %s", e.bracketType)
}

// == Date/Time Errors == //

type InvalidDateFormatError struct {
    format string
    err    error
}

func (e *InvalidDateFormatError) Error() string {
    return fmt.Sprintf("invalid date format: %s, failed to parse with error: %v", e.format, e.err)
}

type InvalidTimeFormatError struct {
    format string
    err    error
}

func (e *InvalidTimeFormatError) Error() string {
    return fmt.Sprintf("invalid time format: %s, failed to parse with error: %v", e.format, e.err)
}

// == Color Errors == //

type InvalidColorError struct {
    color Color
}

func (e *InvalidColorError) Error() string {
    return fmt.Sprintf("invalid color: %v", e.color)
}

// == File Errors == //

var FileNotSpecifiedError = errors.New("filename not provided to NewFileLogger")

type FileNotFoundError struct {
    filename string
}

func (e *FileNotFoundError) Error() string {
    return fmt.Sprintf("file not found for FileLogger: %s", e.filename)
}
