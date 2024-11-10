package ultralogger

type Color interface {
    Colorize(str string) string
}
