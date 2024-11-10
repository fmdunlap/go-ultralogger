package ultralogger

import "fmt"

var ansiReset = []byte("\033[0m")

var ansiCSInit = []byte("\033[")
var ansiCSEnd = byte('m')
var ansiCSSeparator = byte(';')

// TODO: 256 color (maybe)

var ColorRed = AnsiColor{code: []byte("31")}
var ColorGreen = AnsiColor{code: []byte("32")}
var ColorYellow = AnsiColor{code: []byte("33")}
var ColorBlue = AnsiColor{code: []byte("34")}
var ColorMagenta = AnsiColor{code: []byte("35")}
var ColorCyan = AnsiColor{code: []byte("36")}
var ColorWhite = AnsiColor{code: []byte("37")}
var ColorDefault = AnsiColor{code: []byte("39")}

type AnsiColor struct {
    code       []byte
    settings   []ansiSetting
    background AnsiBackgroundColor
}

func RGB(r, g, b int) AnsiColor {
    return AnsiColor{
        code:     []byte(fmt.Sprintf("38;2;%d;%d;%d", r, g, b)),
        settings: []ansiSetting{},
    }
}

func (ac AnsiColor) SetBackground(background AnsiBackgroundColor) AnsiColor {
    ac.background = background
    return ac
}

func (ac AnsiColor) Bold() AnsiColor {
    return AnsiColor{
        code:       ac.code,
        settings:   append(ac.settings, AnsiBold),
        background: ac.background,
    }
}

func (ac AnsiColor) Dim() AnsiColor {
    return AnsiColor{
        code:       ac.code,
        settings:   append(ac.settings, AnsiDim),
        background: ac.background,
    }
}

func (ac AnsiColor) Italic() AnsiColor {
    return AnsiColor{
        code:       ac.code,
        settings:   append(ac.settings, AnsiItalic),
        background: ac.background,
    }
}

func (ac AnsiColor) Underline() AnsiColor {
    return AnsiColor{
        code:       ac.code,
        settings:   append(ac.settings, AnsiUnderline),
        background: ac.background,
    }
}

func (ac AnsiColor) SlowBlink() AnsiColor {
    return AnsiColor{
        code:       ac.code,
        settings:   append(ac.settings, AnsiSlowBlink),
        background: ac.background,
    }
}

// TODO: Benchmark different ways of doing this.
// Went for the single buffer approach for now.
func (ac AnsiColor) Colorize(str string) string {
    buf := make([]byte, ac.totalBufferLength(str))
    cursor := 0

    copy(buf, ansiCSInit)
    cursor += len(ansiCSInit)

    for _, setting := range ac.settings {
        copy(buf[cursor:], setting)
        cursor += len(setting)
        buf[cursor] = ansiCSSeparator
        cursor++
    }

    if ac.background != nil {
        copy(buf[cursor:], ac.background)
        cursor += len(ac.background)
        buf[cursor] = ansiCSSeparator
        cursor++
    }

    copy(buf[cursor:], ac.code)
    cursor += len(ac.code)
    buf[cursor] = ansiCSEnd
    cursor++

    copy(buf[cursor:], str)
    cursor += len(str)

    copy(buf[cursor:], ansiReset)
    cursor += len(ansiReset)

    return string(buf)
}

func (ac AnsiColor) totalBufferLength(input string) int {
    settingsLength := 0
    for _, setting := range ac.settings {
        settingsLength += len(setting) + 1
    }
    backgroundLength := 0
    if ac.background != nil {
        backgroundLength = len(ac.background) + 1
    }

    return len(ansiCSInit) + settingsLength + backgroundLength + len(ac.code) + 1 + len(input) + len(ansiReset)
}
