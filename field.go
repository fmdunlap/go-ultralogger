package ultralogger

type Field interface {
    FieldFormatter() (FieldFormatter, error)
}

type FieldResult struct {
    Name string
    Data any
}

type FieldFormatter func(mCtx LogLineContext, outputFormat OutputFormat, data any) *FieldResult
