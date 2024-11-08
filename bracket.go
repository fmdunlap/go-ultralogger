package ultralogger

// BracketType is a type representing the type of bracket to use for a log message.
//
// It can be one of the following:
//   - BracketTypeNone
//   - BracketTypeSquare
//   - BracketTypeRound
//   - BracketTypeCurly
//   - BracketTypeAngle
//
// BracketTypeNone - "test"
// BracketTypeSquare - "[test]"
// BracketTypeRound - "(test)"
// BracketTypeCurly - "{test}"
// BracketTypeAngle - "<test>"
//
// TODO: Make brackets more flexible via a BracketType interface and funcs defined on the BracketType for open, close,
// etc.
type BracketType int

const (
    BracketTypeNone BracketType = iota
    BracketTypeSquare
    BracketTypeRound
    BracketTypeCurly
    BracketTypeAngle
)

func (b BracketType) String() string {
    switch b {
    case BracketTypeNone:
        return ""
    case BracketTypeSquare:
        return "[]"
    case BracketTypeRound:
        return "()"
    case BracketTypeCurly:
        return "{}"
    case BracketTypeAngle:
        return "<>"
    default:
        return ""
    }
}

// Open returns the opening bracket for the given BracketType.
func (b BracketType) Open() string {
    switch b {
    case BracketTypeNone:
        return ""
    case BracketTypeSquare:
        return "["
    case BracketTypeRound:
        return "("
    case BracketTypeCurly:
        return "{"
    case BracketTypeAngle:
        return "<"
    default:
        return ""
    }
}

// Close returns the closing bracket for the given BracketType.
func (b BracketType) Close() string {
    switch b {
    case BracketTypeNone:
        return ""
    case BracketTypeSquare:
        return "]"
    case BracketTypeRound:
        return ")"
    case BracketTypeCurly:
        return "}"
    case BracketTypeAngle:
        return ">"
    default:
        return ""
    }
}

// wrap returns the string representation of a log message with the given BracketType.
func (b BracketType) wrap(content string) string {
    return b.Open() + content + b.Close()
}

// validBracketType checks if the provided BracketType is valid.
func validBracketType(bracketType BracketType) bool {
    return int(bracketType) >= 0 && int(bracketType) <= 4
}
