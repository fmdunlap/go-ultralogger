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

var FileNotSpecifiedError = errors.New("filename not provided to NewFileLogger")

type FileNotFoundError struct {
    filename string
}

func (e *FileNotFoundError) Error() string {
    return fmt.Sprintf("file not found for FileLogger: %s", e.filename)
}

var ColorizationNotSupportedError = errors.New("formatter does not support colorization")

type LevelParsingError struct {
    level string
}

func (e *LevelParsingError) Error() string {
    return fmt.Sprintf("invalid level: %s", e.level)
}
