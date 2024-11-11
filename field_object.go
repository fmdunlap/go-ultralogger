package ultralogger

type FieldObect[T any] struct {
    name   string
    format FieldFormatter
}

func NewObjectField[T any](name string, formatter ObjectFieldFormatter[T]) Field {
    return &FieldObect[T]{
        name: name,
        format: func(mCtx LogLineContext, outputFormat OutputFormat, data any) *FieldResult {
            data, ok := data.(T)
            if !ok {
                return nil
            }
            return formatter(mCtx, outputFormat, data.(T))
        },
    }
}

func (f *FieldObect[T]) FieldFormatter() (FieldFormatter, error) {
    return f.format, nil
}

type ObjectFieldFormatter[T any] func(mCtx LogLineContext, outputFormat OutputFormat, data T) *FieldResult
