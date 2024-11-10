package field

func NewDateTimeField(dateTimeFormat string) *DateTimeField {
    dtf := &DateTimeField{
        dateTimeFormat: dateTimeFormat,
        clock:          &realClock{},
    }

    return dtf
}

type DateTimeField struct {
    dateTimeFormat string
    clock          clock
}

func (f *DateTimeField) FieldPrinter() (FieldPrinterFunc, error) {
    // TODO: Make a check here for invalid date time format strings.

    return func(info PrintArgs) string {
        return f.clock.Now().Format(f.dateTimeFormat)
    }, nil
}

func (f *DateTimeField) SetDateTimeFormat(format string) Field {
    f.dateTimeFormat = format
    return f
}
