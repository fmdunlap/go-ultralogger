package ultralogger

import (
    "errors"
    "fmt"
)

type LoggerInitializationError struct {
    err error
}

func (e *LoggerInitializationError) Error() string {
    return fmt.Sprintf("error initializing logger: %v", e.err)
}

func (e *LoggerInitializationError) Unwrap() error {
    return e.err
}

var FileNotSpecifiedError = errors.New("filename not provided to NewFileLogger")

type FileNotFoundError struct {
    filename string
}

func (e *FileNotFoundError) Error() string {
    return fmt.Sprintf("file not found for FileLogger: %s", e.filename)
}

var ColorizationNotSupportedError = errors.New("formatter does not support colorization")

type MissingLevelColorError struct {
    level Level
}

func (e *MissingLevelColorError) Error() string {
    return fmt.Sprintf("missing color for level: %v", e.level)
}

type LevelParsingError struct {
    level string
}

func (e *LevelParsingError) Error() string {
    return fmt.Sprintf("invalid level: %s", e.level)
}

type FieldFormatterError struct {
    field Field
    err   error
}

func (e *FieldFormatterError) Error() string {
    return fmt.Sprintf("error formatting field: %v, err=%v", e.field, e.err)
}

func (e *FieldFormatterError) Unwrap() error {
    return e.err
}

type InvalidOutputFormatError struct {
    outputFormat OutputFormat
}

func (e *InvalidOutputFormatError) Error() string {
    return fmt.Sprintf("invalid output format: %v", e.outputFormat)
}

type AmbiguousDestinationError struct{}

func (e *AmbiguousDestinationError) Error() string {
    return "formatters have ambiguous destinations"
}

type InvalidFieldDataTypeError struct {
    field string
}

func (e *InvalidFieldDataTypeError) Error() string {
    return fmt.Sprintf("invalid field data for field: %v", e.field)
}

var EmptyFieldNameError = errors.New("field name cannot be empty")

var NilFormatterError = errors.New("formatter cannot be nil")
