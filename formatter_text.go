package ultralogger

import "fmt"

type TextFormatter struct {
    Fields         []Field
    FieldSeparator string
}

func (f *TextFormatter) FormatLogLine(mCtx LogLineContext, data any) ([]byte, error) {
    lineData := make([]byte, 0)

    for i, field := range f.Fields {
        fieldFormatter, err := field.FieldFormatter()
        if err != nil {
            return nil, &FieldFormatterError{field: field, err: err}
        }

        res := fieldFormatter(mCtx, OutputFormatText, data)

        if res == nil {
            continue
        }

        if i < len(f.Fields)-1 {
            lineData = fmt.Append(lineData, fieldFormatter(mCtx, OutputFormatText, data).Data, " ")
        } else {
            lineData = fmt.Append(lineData, fieldFormatter(mCtx, OutputFormatText, data).Data)
        }
    }

    return lineData, nil
}
