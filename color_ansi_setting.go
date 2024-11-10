package ultralogger

var AnsiBold = ansiSetting("1")
var AnsiDim = ansiSetting("2")
var AnsiItalic = ansiSetting("3")
var AnsiUnderline = ansiSetting("4")
var AnsiSlowBlink = ansiSetting("5")
var AnsiStrikethrough = ansiSetting("9")

type ansiSetting = []byte
