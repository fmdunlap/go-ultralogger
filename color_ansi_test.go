package ultralogger

import (
    "fmt"
    "testing"
)

func TestAnsiColor_Colorize(t *testing.T) {
    tests := []struct {
        name string
        msg  string
        c    AnsiColor
        want string
    }{
        {
            name: "ColorRed",
            msg:  "test",
            c:    ColorRed,
            want: "\033[31mtest\033[0m",
        },
        {
            name: "Bold",
            msg:  "test",
            c: AnsiColor{
                code:     []byte("31"),
                settings: []ansiSetting{AnsiBold},
            },
            want: "\033[1;31mtest\033[0m",
        },
        {
            name: "Dim",
            msg:  "test",
            c: AnsiColor{
                code:     []byte("31"),
                settings: []ansiSetting{AnsiDim},
            },
            want: "\033[2;31mtest\033[0m",
        },
        {
            name: "Italic",
            msg:  "test",
            c: AnsiColor{
                code:     []byte("31"),
                settings: []ansiSetting{AnsiItalic},
            },
            want: "\033[3;31mtest\033[0m",
        },
        {
            name: "Underline",
            msg:  "test",
            c: AnsiColor{
                code:     []byte("31"),
                settings: []ansiSetting{AnsiUnderline},
            },
            want: "\033[4;31mtest\033[0m",
        },
        {
            name: "SlowBlink",
            msg:  "test",
            c: AnsiColor{
                code:     []byte("31"),
                settings: []ansiSetting{AnsiSlowBlink},
            },
            want: "\033[5;31mtest\033[0m",
        },
        {
            name: "Strikethrough",
            msg:  "test",
            c: AnsiColor{
                code:     []byte("31"),
                settings: []ansiSetting{AnsiStrikethrough},
            },
            want: "\033[9;31mtest\033[0m",
        },
        {
            name: "Multiple Settings",
            msg:  "test",
            c: AnsiColor{
                code:     []byte("31"),
                settings: []ansiSetting{AnsiBold, AnsiItalic, AnsiUnderline, AnsiSlowBlink, AnsiStrikethrough},
            },
            want: "\033[1;3;4;5;9;31mtest\033[0m",
        },
        {
            name: "RGB",
            msg:  "test",
            c:    RGB(138, 206, 0),
            want: "\033[38;2;138;206;0mtest\033[0m",
        },
        {
            name: "BackgroundRed",
            msg:  "test",
            c: AnsiColor{
                code:       []byte("30"),
                settings:   []ansiSetting{},
                background: BackgroundRed,
            },
            want: "\033[41;30mtest\033[0m",
        },
        {
            name: "BackgroundRGB",
            msg:  "test",
            c: AnsiColor{
                code:       []byte("30"),
                settings:   []ansiSetting{},
                background: BackgroundRGB(138, 206, 0),
            },
            want: "\033[48;2;138;206;0;30mtest\033[0m",
        },
        {
            name: "BackgroundRGB + Bold",
            msg:  "test",
            c: AnsiColor{
                code:       []byte("30"),
                settings:   []ansiSetting{AnsiBold},
                background: BackgroundRGB(138, 206, 0),
            },
            want: "\033[1;48;2;138;206;0;30mtest\033[0m",
        },
        {
            name: "BackgroundRed + Multiple Settings",
            msg:  "test",
            c: AnsiColor{
                code:       []byte("30"),
                settings:   []ansiSetting{AnsiBold, AnsiItalic, AnsiUnderline, AnsiSlowBlink, AnsiStrikethrough},
                background: BackgroundRed,
            },
            want: "\033[1;3;4;5;9;41;30mtest\033[0m",
        },
        {
            name: "BackgroundRGB + Multiple Settings",
            msg:  "test",
            c: AnsiColor{
                code:       []byte("30"),
                settings:   []ansiSetting{AnsiBold, AnsiItalic, AnsiUnderline, AnsiSlowBlink, AnsiStrikethrough},
                background: BackgroundRGB(138, 206, 0),
            },
            want: "\033[1;3;4;5;9;48;2;138;206;0;30mtest\033[0m",
        },
        {
            name: "ColorRGB + BackgroundRGB",
            msg:  "test",
            c:    RGB(138, 206, 0).SetBackground(BackgroundRGB(255, 0, 0)),
            want: "\033[48;2;255;0;0;38;2;138;206;0mtest\033[0m",
        },
        {
            name: "ColorRGB + BackgroundRGB + Bold",
            msg:  "test",
            c:    RGB(138, 206, 0).SetBackground(BackgroundRGB(255, 0, 0)).Bold(),
            want: "\033[1;48;2;255;0;0;38;2;138;206;0mtest\033[0m",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            fmt.Println(tt.c.Colorize(tt.msg))
            if got := tt.c.Colorize(tt.msg); got != tt.want {
                fmt.Println("Got:  ", []byte(got))
                fmt.Println("Want: ", []byte(tt.want))
                t.Errorf("Colorize() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestAnsiColor_totalBufferLength(t *testing.T) {
    tests := []struct {
        name  string
        c     AnsiColor
        input string
        want  int
    }{
        {
            name: "No Settings",
            c: AnsiColor{
                code:       []byte("31"),
                settings:   []ansiSetting{},
                background: nil,
                // output:     "\033[31mtest\033[0m",
            },
            input: "test",
            want:  13,
        },
        {
            name: "Bold",
            c: AnsiColor{
                code:       []byte("31"),
                settings:   []ansiSetting{AnsiBold},
                background: nil,
                // output:     "\033[1;31mtest\033[0m",
            },
            input: "test",
            want:  15,
        },
        {
            name: "Multiple Settings",
            c: AnsiColor{
                code:       []byte("31"),
                settings:   []ansiSetting{AnsiBold, AnsiItalic, AnsiUnderline, AnsiSlowBlink, AnsiStrikethrough},
                background: nil,
                // output:     "\033[1;3;4;5;9;31mtest\033[0m",
            },
            input: "test",
            want:  23,
        },
        {
            name: "RGB",
            c:    RGB(138, 206, 0),
            // output: "\033[38;2;138;206;0mtest\033[0m",
            input: "test",
            want:  25,
        },
        {
            name: "BackgroundRed",
            c: AnsiColor{
                code:       []byte("30"),
                settings:   []ansiSetting{},
                background: BackgroundRed,
                // output:     "\033[41;30mtest\033[0m",
            },
            input: "test",
            want:  16,
        },
        {
            name: "BackgroundRGB",
            c: AnsiColor{
                code:       []byte("30"),
                settings:   []ansiSetting{},
                background: BackgroundRGB(138, 206, 0),
                // output:     "\033[48;2;138;206;0;30mtest\033[0m",
            },
            input: "test",
            want:  28,
        },
        {
            name:  "ColorRGB + BackgroundRGB",
            c:     RGB(138, 206, 0).SetBackground(BackgroundRGB(255, 0, 0)),
            input: "test",
            want:  38,
        },
        {
            name:  "ColorRGB + BackgroundRGB + Bold",
            c:     RGB(138, 206, 0).SetBackground(BackgroundRGB(255, 0, 0)).Bold(),
            input: "test",
            want:  40,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.c.totalBufferLength(tt.input); got != tt.want {
                t.Errorf("totalBufferLength() = %v, want %v", got, tt.want)
            }
        })
    }
}
