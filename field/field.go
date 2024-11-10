package field

import (
    "github.com/fmdunlap/go-ultralogger/level"
)

type FieldPrinterFunc func(PrintArgs) string

type PrintArgs struct {
    Level level.Level
}

type Field interface {
    FieldPrinter() (FieldPrinterFunc, error)
}
