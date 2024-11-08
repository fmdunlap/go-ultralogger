package ultralogger

// Color is a type representing an ANSI terminal color for a log message.
// It can be one of the following:
//   - ColorRed
//   - ColorGreen
//   - ColorYellow
//   - ColorBlue
//   - ColorMagenta
//   - ColorCyan
//   - ColorGray
//   - ColorWhite
//   - ColorDefault
//
// ColorRed - "\033[31m"
// ColorGreen - "\033[32m"
// ColorYellow - "\033[33m"
// ColorBlue - "\033[34m"
// ColorMagenta - "\033[35m"
// ColorCyan - "\033[36m"
// ColorGray - "\033[2;37m"
// ColorWhite - "\033[37m"
// ColorDefault - "\033[39m"
//
// TODO: There's lots of room for improvement here. Background colors, bright/dim, readability, etc. Open a PR!
type Color int

const (
    ColorRed = iota
    ColorGreen
    ColorYellow
    ColorBlue
    ColorMagenta
    ColorCyan
    ColorGray
    ColorWhite
    ColorDefault
)

var reset = []byte("\033[0m")

// code returns the escape sequence for the given Color.
func (c Color) code() []byte {
    switch c {
    case ColorRed:
        return []byte("\033[31m")
    case ColorGreen:
        return []byte("\033[32m")
    case ColorYellow:
        return []byte("\033[33m")
    case ColorBlue:
        return []byte("\033[34m")
    case ColorMagenta:
        return []byte("\033[35m")
    case ColorCyan:
        return []byte("\033[36m")
    case ColorGray:
        return []byte("\033[2;37m")
    case ColorWhite:
        return []byte("\033[37m")
    case ColorDefault:
        return []byte("\033[39m")
    default:
        return []byte("\033[39m")
    }
}

func (c Color) String() string {
    switch c {
    case ColorRed:
        return "RED"
    case ColorGreen:
        return "GREEN"
    case ColorYellow:
        return "YELLOW"
    case ColorBlue:
        return "BLUE"
    case ColorMagenta:
        return "MAGENTA"
    case ColorCyan:
        return "CYAN"
    case ColorGray:
        return "GRAY"
    case ColorWhite:
        return "WHITE"
    case ColorDefault:
        return "DEFAULT"
    default:
        return "UNKNOWN"
    }
}

// colorize returns the string representation of a log message with the given color.
func colorize(color Color, msg string) string {
    return string(color.code()) + msg + string(reset)
}

// validColor checks if the provided color is a known color.
func validColor(color Color) bool {
    return int(color) >= 0 && color <= ColorDefault
}
