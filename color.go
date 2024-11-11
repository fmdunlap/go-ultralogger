package ultralogger

type Color interface {
    Colorize(str []byte) []byte
}
