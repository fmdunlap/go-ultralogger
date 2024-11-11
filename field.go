package ultralogger

import (
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// Field A log field is a piece of data that is added to a log line. It can be a simple value, or an object.
// The FieldFormatter is responsible for formatting the data into a FieldResult.
type Field interface {
    FieldFormatter() (FieldFormatter, error)
}

// FieldResult is the result of formatting a field. It is used by the logger to annotate entries with metadata, suhc as
// the name of the field.
type FieldResult struct {
    Name string
    Data any
}

// FieldFormatter is a function that formats a field. It takes a LogLineArgs and the data to be formatted, and returns
// a FieldResult.
type FieldFormatter func(
    args LogLineArgs,
    data any,
) (FieldResult, error)

type ObjectField[T any] struct {
    format FieldFormatter
}

type ObjectFieldFormatter[T any] func(
    args LogLineArgs,
    data T,
) any

func (f ObjectField[T]) FieldFormatter() (FieldFormatter, error) {
    return f.format, nil
}

func NewObjectField[T any](name string, formatter ObjectFieldFormatter[T]) (ObjectField[T], error) {
    if name == "" {
        return ObjectField[T]{}, EmptyFieldNameError
    }
    if formatter == nil {
        return ObjectField[T]{}, NilFormatterError
    }
    return ObjectField[T]{
        format: func(args LogLineArgs, data any) (FieldResult, error) {
            result := FieldResult{
                Name: name,
            }

            _, ok := data.(T)
            if !ok {
                return result, &InvalidFieldDataTypeError{
                    field: name,
                }
            }
            result.Data = formatter(args, data.(T))

            return result, nil
        },
    }, nil
}

func NewStringField(name string) (Field, error) {
    return NewObjectField[string](
        name,
        func(args LogLineArgs, data string) any {
            return data
        },
    )
}

func NewBoolField(name string) (Field, error) {
    return NewObjectField[bool](
        name,
        func(args LogLineArgs, data bool) any {
            if args.OutputFormat == OutputFormatText {
                if data {
                    return "true"
                }
                return "false"
            }
            return data
        },
    )
}

func NewTimeField(name, format string) (Field, error) {
    return NewObjectField[time.Time](
        name,
        func(args LogLineArgs, data time.Time) any {
            if args.OutputFormat == OutputFormatText {
                return data.Format(format)
            }
            return data
        },
    )
}

func NewIntField(name string) (Field, error) {
    return NewObjectField[int](
        name,
        func(args LogLineArgs, data int) any {
            if args.OutputFormat == OutputFormatText {
                return strconv.Itoa(data)
            }
            return data
        },
    )
}

func NewFloatField(name string) (Field, error) {
    return NewObjectField[float64](
        name,
        func(args LogLineArgs, data float64) any {
            if args.OutputFormat == OutputFormatText {
                return strconv.FormatFloat(data, 'f', -1, 64)
            }
            return data
        },
    )
}

func NewDurationField(name string) (Field, error) {
    return NewObjectField[time.Duration](
        name,
        func(args LogLineArgs, data time.Duration) any {
            if args.OutputFormat == OutputFormatText {
                return data.String()
            }
            return data
        },
    )
}

func NewErrorField(name string) (Field, error) {
    return NewObjectField[error](
        name,
        func(args LogLineArgs, data error) any {
            if args.OutputFormat == OutputFormatText {
                return data.Error()
            }
            return data
        },
    )
}

func NewArrayField[T any](name string, formatter ObjectFieldFormatter[T]) (Field, error) {
    if name == "" {
        return ObjectField[[]T]{}, EmptyFieldNameError
    }
    return NewObjectField[[]T](
        name,
        func(args LogLineArgs, data []T) any {
            res := make([]any, len(data))
            for i, v := range data {
                res[i] = formatter(args, v)
            }

            if args.OutputFormat == OutputFormatText {
                b := strings.Builder{}
                b.WriteString("[")
                for i, v := range res {
                    b.WriteString(fmt.Sprintf("%v", v))
                    if i < len(data)-1 {
                        b.WriteString(", ")
                    }
                }
                b.WriteString("]")
                return b.String()
            }

            fmt.Println("res", res)

            return res
        },
    )
}

func NewMapField[K comparable, V any](name string, keyFormatter ObjectFieldFormatter[K], valueFormatter ObjectFieldFormatter[V]) (Field, error) {
    if name == "" {
        return ObjectField[map[K]V]{}, EmptyFieldNameError
    }
    if keyFormatter == nil {
        return ObjectField[map[K]V]{}, NilFormatterError
    }
    if valueFormatter == nil {
        return ObjectField[map[K]V]{}, NilFormatterError
    }
    return NewObjectField[map[K]V](
        name,
        func(args LogLineArgs, data map[K]V) any {
            res := make(map[any]any)
            for k, v := range data {
                res[keyFormatter(args, k)] = valueFormatter(args, v)
            }

            if args.OutputFormat != OutputFormatText {
                validMap := make(map[string]any)
                for k, v := range res {
                    validMap[fmt.Sprintf("%v", k)] = v
                }
                return validMap
            }

            return res
        },
    )

}

type CurrentTimeField struct {
    name      string
    fmtString string
    clock     clock
}

func (f *CurrentTimeField) FieldFormatter() (FieldFormatter, error) {
    return f.format, nil
}

func (f *CurrentTimeField) format(args LogLineArgs, _ any) (FieldResult, error) {
    result := FieldResult{
        Name: f.name,
    }

    now := f.clock.Now()

    switch args.OutputFormat {
    case OutputFormatJSON:
        result.Data = now
    case OutputFormatText:
        result.Data = now.Format(f.fmtString)
    }

    return result, nil
}

func NewCurrentTimeField(name, format string) (Field, error) {
    if name == "" {
        return &CurrentTimeField{}, EmptyFieldNameError
    }

    ctf := &CurrentTimeField{
        name:      name,
        fmtString: format,
        clock:     &realClock{},
    }

    return ctf, nil
}

func NewLevelField(bracket Bracket) *FieldLevel {
    return &FieldLevel{
        bracket: bracket,
    }
}

type FieldLevel struct {
    bracket      Bracket
    levelStrings map[Level]string
}

func (f *FieldLevel) FieldFormatter() (FieldFormatter, error) {
    if f.levelStrings == nil {
        f.levelStrings = make(map[Level]string)

        for _, lvl := range AllLevels() {
            f.levelStrings[lvl] = f.bracket.Wrap(lvl.String())
        }
    }

    return f.format, nil
}

func (f *FieldLevel) format(args LogLineArgs, _ any) (FieldResult, error) {
    return FieldResult{
        Name: "Level",
        Data: f.levelStrings[args.Level],
    }, nil
}

type FieldMessage struct{}

func (f *FieldMessage) FieldFormatter() (FieldFormatter, error) {
    return f.format, nil
}

func (f *FieldMessage) format(_ LogLineArgs, message any) (FieldResult, error) {
    result := FieldResult{
        Name: "Message",
    }

    switch message.(type) {
    case string:
        result.Data = message.(string)
    case fmt.Stringer:
        result.Data = message.(fmt.Stringer).String()
    default:
        return result, &InvalidFieldDataTypeError{
            field: "Message",
        }
    }

    return result, nil
}

// TODO: There's definitely more to be added to Request & Response logging.

type RequestFieldSettings struct {
    timeFormat string

    LogReceivedAt bool
    LogMethod     bool
    LogPath       bool
    LogSourceIP   bool
}

type RequestLogEntry struct {
    ReceivedAt time.Time
    Method     string
    Path       string
    SourceIP   string
}

func (r *RequestLogEntry) String(timeFmt string) string {
    parts := []string{}
    if !r.ReceivedAt.IsZero() {
        parts = append(parts, r.ReceivedAt.Format(timeFmt))
    }
    if r.Method != "" {
        parts = append(parts, r.Method)
    }
    if r.Path != "" {
        parts = append(parts, r.Path)
    }
    if r.SourceIP != "" {
        parts = append(parts, r.SourceIP)
    }
    return strings.Join(parts, " ")
}

func NewRequestField(name string, settings RequestFieldSettings) (Field, error) {
    return NewObjectField[*http.Request](
        name,
        func(args LogLineArgs, data *http.Request) any {
            logEntry := RequestLogEntry{}

            if settings.LogReceivedAt {
                logEntry.ReceivedAt = time.Now()
            }

            if settings.LogSourceIP {
                logEntry.SourceIP = data.RemoteAddr
            }

            if settings.LogMethod {
                logEntry.Method = data.Method
            }

            if settings.LogPath {
                logEntry.Path = data.URL.Path
            }

            if args.OutputFormat == OutputFormatText {
                return logEntry.String(settings.timeFormat)
            }
            return logEntry
        },
    )
}

type ResponseLogSettings struct {
    LogStatus bool
    LogPath   bool
}

type ResponseLogEntry struct {
    StatusCode int
    Status     string
    Path       string
}

func (r *ResponseLogEntry) String() string {
    parts := []string{}
    if r.StatusCode != 0 {
        parts = append(parts, strconv.Itoa(r.StatusCode))
    }
    return strings.Join(parts, " ")
}

func NewResponseField(name string, settings ResponseLogSettings) (Field, error) {
    return NewObjectField[*http.Response](
        name,
        func(args LogLineArgs, data *http.Response) any {
            logEntry := ResponseLogEntry{}

            if settings.LogStatus {
                logEntry.StatusCode = data.StatusCode
                logEntry.Status = data.Status
            }

            if settings.LogPath {
                logEntry.Path = data.Request.URL.Path
            }

            if args.OutputFormat == OutputFormatText {
                return logEntry.String()
            }
            return logEntry
        },
    )
}
