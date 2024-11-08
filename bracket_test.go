package ultralogger

import "testing"

func TestBracketType_Close(t *testing.T) {
    tests := []struct {
        name string
        b    BracketType
        want string
    }{
        {"BracketTypeNone", BracketTypeNone, ""},
        {"BracketTypeSquare", BracketTypeSquare, "]"},
        {"BracketTypeRound", BracketTypeRound, ")"},
        {"BracketTypeCurly", BracketTypeCurly, "}"},
        {"BracketTypeAngle", BracketTypeAngle, ">"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.b.Close(); got != tt.want {
                t.Errorf("Close() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestBracketType_Open(t *testing.T) {
    tests := []struct {
        name string
        b    BracketType
        want string
    }{
        {"BracketTypeNone", BracketTypeNone, ""},
        {"BracketTypeSquare", BracketTypeSquare, "["},
        {"BracketTypeRound", BracketTypeRound, "("},
        {"BracketTypeCurly", BracketTypeCurly, "{"},
        {"BracketTypeAngle", BracketTypeAngle, "<"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.b.Open(); got != tt.want {
                t.Errorf("Open() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestBracketType_String(t *testing.T) {
    tests := []struct {
        name string
        b    BracketType
        want string
    }{
        {"BracketTypeNone", BracketTypeNone, ""},
        {"BracketTypeSquare", BracketTypeSquare, "[]"},
        {"BracketTypeRound", BracketTypeRound, "()"},
        {"BracketTypeCurly", BracketTypeCurly, "{}"},
        {"BracketTypeAngle", BracketTypeAngle, "<>"},
        {"BracketTypeUnknown", BracketType(42), ""},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.b.String(); got != tt.want {
                t.Errorf("String() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestBracketType_Wrap(t *testing.T) {
    tests := []struct {
        name    string
        b       BracketType
        content string
        want    string
    }{
        {"BracketTypeNone", BracketTypeNone, "test", "test"},
        {"BracketTypeSquare", BracketTypeSquare, "test", "[test]"},
        {"BracketTypeRound", BracketTypeRound, "test", "(test)"},
        {"BracketTypeCurly", BracketTypeCurly, "test", "{test}"},
        {"BracketTypeAngle", BracketTypeAngle, "test", "<test>"},
        {"BracketTypeUnknown", BracketType(42), "test", "test"},
        {"EmptyContent-BracketTypeNone", BracketTypeNone, "", ""},
        {"EmptyContent-BracketTypeSquare", BracketTypeSquare, "", "[]"},
        {"EmptyContent-BracketTypeRound", BracketTypeRound, "", "()"},
        {"EmptyContent-BracketTypeCurly", BracketTypeCurly, "", "{}"},
        {"EmptyContent-BracketTypeAngle", BracketTypeAngle, "", "<>"},
        {"EmptyContent-BracketTypeUnknown", BracketType(42), "", ""},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.b.wrap(tt.content); got != tt.want {
                t.Errorf("wrap() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_validBracketType(t *testing.T) {
    tests := []struct {
        name        string
        bracketType BracketType
        want        bool
    }{
        {"BracketTypeNone", BracketTypeNone, true},
        {"BracketTypeSquare", BracketTypeSquare, true},
        {"BracketTypeRound", BracketTypeRound, true},
        {"BracketTypeCurly", BracketTypeCurly, true},
        {"BracketTypeAngle", BracketTypeAngle, true},
        {"InvalidBracketTypeNonExistent", BracketType(10), false},
        {"InvalidBracketTypeNegative", BracketType(-1), false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := validBracketType(tt.bracketType); got != tt.want {
                t.Errorf("validBracketType() = %v, want %v", got, tt.want)
            }
        })
    }
}
