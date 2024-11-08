package ultralogger

import "time"

type clock interface {
    Now() time.Time
}

type realClock struct{}

func (c *realClock) Now() time.Time {
    return time.Now()
}
