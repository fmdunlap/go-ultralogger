# UltraLogger

[![Logo](/git/logo.png)](https://github.com/forrestmdunlap/ultralogger)

## Overview
UltraLogger is a versatile, flexible, and efficient logging library for Go that supports multiple log levels, custom formatting, and output destinations. It is designed to be easy to use while providing advanced features such as colorized output, customizable log levels, and detailed error handling.

## Features
- **Multiple Log Levels**: Supports DEBUG, INFO, WARN, ERROR, PANIC, and FATAL levels.
- **Custom Formatting**: Flexible formatting options including custom date/time formats, tag padding, and bracket types.
- **Colorization**: Colorizes log output based on severity level for better visibility.
- **Output Redirection**: Supports writing logs to various io.Writer destinations such as files or standard output.
- **Silent Mode**: Allows disabling all logging when set.
- **Error Handling**: Implements robust error handling and fallback mechanisms for logging errors.

## Installation
To install UltraLogger, use `go get`:
```sh
go get github.com/forrestmdunlap/ultralogger
```

## Usage
Here's a basic example of how to use UltraLogger:

```go
package main

import (
    "os"
    "github.com/forrestmdunlap/ultralogger"
)

func main() {
    logger := ultralogger.NewUltraLogger(os.Stdout)
    logger.Info("This is an info message.")
    logger.Debug("This is a debug message.")
    logger.Warn("This is a warning message.")
}
```

TODO: Readme :)