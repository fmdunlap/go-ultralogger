package ultralogger

import "fmt"

type FieldMessage struct{}

func (f *FieldMessage) FieldFormatter() (FieldFormatter, error) {
    return f.format, nil
}

func (f *FieldMessage) format(_ LogLineContext, outputFormat OutputFormat, message any) *FieldResult {
    result := &FieldResult{
        Name: "Message",
    }

    switch message.(type) {
    case string:
        result.Data = message.(string)
    case fmt.Stringer:
        result.Data = message.(fmt.Stringer).String()
    default:
        return nil
    }

    return result
}
