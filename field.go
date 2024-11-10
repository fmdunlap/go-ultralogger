package ultralogger

type FieldPrinterFunc func(PrintArgs) string

type PrintArgs struct {
    Level Level
}

type Field interface {
    FieldPrinter() (FieldPrinterFunc, error)
}
