package ultralogger

func NewDateTimeField(dateTimeFormat string) *FieldDateTime {
    dtf := &FieldDateTime{
        dateTimeFormat: dateTimeFormat,
        clock:          &realClock{},
    }

    return dtf
}

type FieldDateTime struct {
    dateTimeFormat string
    clock          clock
}

func (f *FieldDateTime) FieldPrinter() (FieldPrinterFunc, error) {
    // TODO: Make a check here for invalid date time format strings.

    return func(info PrintArgs) string {
        return f.clock.Now().Format(f.dateTimeFormat)
    }, nil
}

func (f *FieldDateTime) SetDateTimeFormat(format string) Field {
    f.dateTimeFormat = format
    return f
}
