package color

var Bold = AnsiBackgroundColor("1")
var Dim = AnsiBackgroundColor("2")
var Italic = AnsiBackgroundColor("3")
var Underline = AnsiBackgroundColor("4")
var SlowBlink = AnsiBackgroundColor("5")
var Strikethrough = AnsiBackgroundColor("9")

type ansiSetting = []byte
