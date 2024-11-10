package color

import "fmt"

type AnsiBackgroundColor = []byte

var BackgroundBlack = AnsiBackgroundColor("40")
var BackgroundRed = AnsiBackgroundColor("41")
var BackgroundGreen = AnsiBackgroundColor("42")
var BackgroundYellow = AnsiBackgroundColor("43")
var BackgroundBlue = AnsiBackgroundColor("44")
var BackgroundMagenta = AnsiBackgroundColor("45")
var BackgroundCyan = AnsiBackgroundColor("46")
var BackgroundWhite = AnsiBackgroundColor("47")

func BackgroundRGB(r, g, b int) AnsiBackgroundColor {
    return AnsiBackgroundColor(fmt.Sprintf("48;2;%d;%d;%d", r, g, b))
}
