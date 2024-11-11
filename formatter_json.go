package ultralogger

import (
    "encoding/json"
    "errors"
)

type JSONFormatter struct {
    Fields                 []Field
    destinationInitialized bool
}

func (f *JSONFormatter) FormatLogLine(args LogLineArgs, data any) FormatResult {
    jsonMap := make(map[string]any)

    args.OutputFormat = OutputFormatJSON

FieldLoop:
    for _, field := range f.Fields {
        fieldFormatter, err := field.FieldFormatter()
        if err != nil {
            return FormatResult{nil, &FieldFormatterError{field: field, err: err}}
        }

        fieldResult, err := fieldFormatter(args, data)

        if err != nil {
            var invalidDataTypeError *InvalidFieldDataTypeError
            if errors.As(err, &invalidDataTypeError) {
                continue FieldLoop
            }
        }

        if fieldResult.Data == nil {
            continue FieldLoop
        }

        jsonMap[fieldResult.Name] = fieldResult.Data
    }

    jBytes, err := json.Marshal(jsonMap)
    if err != nil {
        return FormatResult{nil, &FieldFormatterError{field: nil, err: err}}
    }

    return FormatResult{jBytes, nil}
}
