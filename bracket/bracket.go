package bracket

var Angle = SimpleBracket{"<", ">"}
var Square = SimpleBracket{"[", "]"}
var Round = SimpleBracket{"(", ")"}
var Curly = SimpleBracket{"{", "}"}
var None = SimpleBracket{"", ""}

type Bracket interface {
    Open() string
    Close() string
    Wrap(content string) string
}

type SimpleBracket struct {
    open  string
    close string
}

func (sb SimpleBracket) Open() string {
    return sb.open
}

func (sb SimpleBracket) Close() string {
    return sb.close
}

func (sb SimpleBracket) Wrap(content string) string {
    return sb.open + content + sb.close
}
