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

func (f *FieldDateTime) FieldFormatter() (FieldFormatter, error) {
    // TODO: Make a check here for invalid date time format strings.

    return f.format, nil
}

func (f *FieldDateTime) SetDateTimeFormat(format string) Field {
    f.dateTimeFormat = format
    return f
}

func (f *FieldDateTime) format(_ LogLineContext, outputFormat OutputFormat, _ any) *FieldResult {
    result := &FieldResult{
        Name: "DateTime",
    }

    now := f.clock.Now()

    switch outputFormat {
    case OutputFormatJSON:
        result.Data = now
    case OutputFormatText:
        result.Data = now.Format(f.dateTimeFormat)
    }

    return result
}
