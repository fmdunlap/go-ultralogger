package level

import "fmt"

type ParsingError struct {
    level string
}

func (e *ParsingError) Error() string {
    return fmt.Sprintf("invalid level: %s", e.level)
}
