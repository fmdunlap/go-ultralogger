package bracket

import "testing"

func TestSimpleBracket_Close(t *testing.T) {
    type fields struct {
        open  string
        close string
    }
    tests := []struct {
        name   string
        fields fields
        want   string
    }{
        {
            name: "Close returns the closing bracket",
            fields: fields{
                open:  "<",
                close: ">",
            },
            want: ">",
        },
        {
            name: "Empty string bracket returns empty string",
            fields: fields{
                open:  "",
                close: "",
            },
            want: "",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sb := SimpleBracket{
                open:  tt.fields.open,
                close: tt.fields.close,
            }
            if got := sb.Close(); got != tt.want {
                t.Errorf("Close() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestSimpleBracket_Open(t *testing.T) {
    type fields struct {
        open  string
        close string
    }
    tests := []struct {
        name   string
        fields fields
        want   string
    }{
        {
            name: "Open returns the opening bracket",
            fields: fields{
                open:  "<",
                close: ">",
            },
            want: "<",
        },
        {
            name: "Empty string bracket returns empty string",
            fields: fields{
                open:  "",
                close: "",
            },
            want: "",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sb := SimpleBracket{
                open:  tt.fields.open,
                close: tt.fields.close,
            }
            if got := sb.Open(); got != tt.want {
                t.Errorf("Open() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestSimpleBracket_Wrap(t *testing.T) {
    type fields struct {
        open  string
        close string
    }
    type args struct {
        content string
    }
    tests := []struct {
        name   string
        fields fields
        args   args
        want   string
    }{
        {
            name: "Content is wrapped",
            fields: fields{
                open:  "<",
                close: ">",
            },
            args: args{
                content: "test",
            },
            want: "<test>",
        },
        {
            name: "Empty content is wrapped",
            fields: fields{
                open:  "<",
                close: ">",
            },
            args: args{
                content: "",
            },
            want: "<>",
        },
        {
            name: "Empty string bracket returns string",
            fields: fields{
                open:  "",
                close: "",
            },
            args: args{
                content: "test",
            },
            want: "test",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sb := SimpleBracket{
                open:  tt.fields.open,
                close: tt.fields.close,
            }
            if got := sb.Wrap(tt.args.content); got != tt.want {
                t.Errorf("Wrap() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestBracketType_BuiltinWrap(t *testing.T) {
    tests := []struct {
        name    string
        b       Bracket
        content string
        want    string
    }{
        {
            name:    "Angle",
            b:       Angle,
            content: "test",
            want:    "<test>",
        },
        {
            name:    "Square",
            b:       Square,
            content: "test",
            want:    "[test]",
        },
        {
            name:    "Round",
            b:       Round,
            content: "test",
            want:    "(test)",
        },
        {
            name:    "Curly",
            b:       Curly,
            content: "test",
            want:    "{test}",
        },
        {
            name:    "None",
            b:       None,
            content: "test",
            want:    "test",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.b.Wrap(tt.content); got != tt.want {
                t.Errorf("Wrap() = %v, want %v", got, tt.want)
            }
        })
    }
}
