# UltraLogger

[![Logo](/git/logo.png)](https://github.com/fmdunlap/ultralogger)

## Overview
UltraLogger is a versatile, flexible, and efficient logging library for Go that supports multiple log levels, custom formatting, and output destinations. It is designed to be easy to use while providing advanced features such as colorized output, customizable log levels, and detailed error handling.

## Features
- **Multiple Log Levels**: DEBUG, INFO, WARN, ERROR, PANIC, and FATAL levels.
- **Custom Formatting**: Flexible formatting options including custom date/time formats, tag padding, and bracket types.
- **Colorization**: Colorizes log output based on severity level for better visibility.
- **Output Redirection**: Supports writing logs to various io.Writer destinations such as files or standard output.
- **Silent Mode**: Allows disabling all logging when set.
- **Error Handling**: Implements robust error handling and fallback mechanisms for logging errors.

## Installation
To install UltraLogger, use `go get`:
```sh
go get github.com/fmdunlap/ultralogger
```

## Usage
Here's a basic example of how to use UltraLogger:

```go
package main

import (
    "os"
    "github.com/fmdunlap/ultralogger"
)

func main() {
    logger := ultralogger.NewUltraLogger(os.Stdout)
    logger.Info("This is an info message.")  // Output: 2006-01-02 15:04:05 <INFO> This is an info message.
    logger.Debug("This is a debug message.") // Output: 2006-01-02 15:04:05 <DEBUG> This is a debug message.
    logger.Warn("This is a warning message.") // Output: 2006-01-02 15:04:05 <WARN> This is a warning message.
}
```

## Configuration

Ultralogger provides various configuration options to customize logging behaviour.

### Minimum Log Level

```go
logger := ultralogger.NewUltraLogger(os.Stdout).MinLogLevel(ultralogger.LogLevelDebug)
```

### Custom Formatting

```go
logger := ultralogger.NewUltraLogger(os.Stdout)

logger.Info("Message") // -> 2006-01-02 15:04:05 <INFO> Message

// Date/Time
logger.SetDateFormat("01/02/2006")
logger.ShowTime(false)

logger.Info("Message") // -> 01/02/2006 <INFO> Message

// Tag
logger.SetTag("MyTag")
logger.TagPadingEnabled(true)
logger.SetTagPadding(10)

logger.Info("Message") // -> 2006-01-02 15:04:05 [MyTag]   <INFO> Message
logger.Error("Error!") // -> 2006-01-02 15:04:05 [MyTag]   <ERROR> Warning!
```

### Arbitrary Output Destinations

```go
// Files
ultralogger.NewFileLogger("somefile.log")

// Stdout
ultralogger.NewStdoutLogger()

// ByteBuffer
buf := new(bytes.Buffer)
ultralogger.NewUltraLogger(buf)
```

### Terminal Colorization

```go
logger := ultralogger.NewUltraLogger(os.Stdout).Colorize(true)

logger.Debug("Debug")
logger.Info("Message")
logger.Warn("Warning")
logger.Error("Error!")
logger.Panic("Panic!")
```



### Silent Mode

TODO: Readme :)