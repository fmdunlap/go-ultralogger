package ultralogger

var BracketAngle = SimpleBracket{"<", ">"}
var BracketSquare = SimpleBracket{"[", "]"}
var BracketRound = SimpleBracket{"(", ")"}
var BracketCurly = SimpleBracket{"{", "}"}
var BracketNone = SimpleBracket{"", ""}

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
