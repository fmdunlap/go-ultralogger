package ultralogger

type OutputFormat string

const (
    OutputFormatJSON OutputFormat = "json"
    OutputFormatText OutputFormat = "text"
)

type LogLineArgs struct {
    Level        Level
    Tag          string
    OutputFormat OutputFormat
}

type FormatResult struct {
    bytes []byte
    err   error
}

type LogLineFormatter interface {
    FormatLogLine(args LogLineArgs, data any) FormatResult
}

type FormatterOption func(f LogLineFormatter) LogLineFormatter

func NewFormatter(outputFormat OutputFormat, fields []Field, opts ...FormatterOption) (LogLineFormatter, error) {
    var f LogLineFormatter

    switch outputFormat {
    case OutputFormatJSON:
        f = &JSONFormatter{Fields: fields}
    case OutputFormatText:
        f = &TextFormatter{Fields: fields}
    default:
        return nil, &InvalidOutputFormatError{outputFormat: outputFormat}
    }

    for _, opt := range opts {
        f = opt(f)
    }

    return f, nil
}

func WithDefaultColorization() FormatterOption {
    return func(f LogLineFormatter) LogLineFormatter {
        return NewColorizedFormatter(f, nil)
    }
}

func WithColorization(colors map[Level]Color) FormatterOption {
    return func(f LogLineFormatter) LogLineFormatter {
        return NewColorizedFormatter(f, colors)
    }
}
