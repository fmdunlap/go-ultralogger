package ultralogger

import "encoding/json"

type JSONFormatter struct {
    Fields                 []Field
    destinationInitialized bool
}

func (f *JSONFormatter) FormatLogLine(mCtx LogLineContext, data any) ([]byte, error) {
    jsonMap := make(map[string]any)

    for _, field := range f.Fields {
        fieldFormatter, err := field.FieldFormatter()
        if err != nil {
            return nil, &FieldFormatterError{field: field, err: err}
        }

        fieldResult := fieldFormatter(mCtx, OutputFormatJSON, data)

        if fieldResult == nil {
            continue
        }

        jsonMap[fieldResult.Name] = fieldResult.Data
    }

    return json.Marshal(jsonMap)
}
