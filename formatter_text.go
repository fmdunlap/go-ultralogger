package ultralogger

import (
    "errors"
    "fmt"
)

type TextFormatter struct {
    Fields         []Field
    FieldSeparator string
}

func (f *TextFormatter) FormatLogLine(args LogLineArgs, data any) FormatResult {
    lineData := make([]byte, 0)

    args.OutputFormat = OutputFormatText

FieldLoop:
    for i, field := range f.Fields {
        fieldFormatter, err := field.FieldFormatter()
        if err != nil {
            return FormatResult{nil, &FieldFormatterError{field: field, err: err}}
        }

        res, err := fieldFormatter(args, data)

        if err != nil {
            var invalidDataTypeError *InvalidFieldDataTypeError
            if errors.As(err, &invalidDataTypeError) {
                continue FieldLoop
            }
        }

        if res.Data == nil {
            continue FieldLoop
        }

        if i < len(f.Fields)-1 {
            lineData = fmt.Append(lineData, res.Data, " ")
        } else {
            lineData = fmt.Append(lineData, res.Data)
        }
    }

    return FormatResult{lineData, nil}
}
