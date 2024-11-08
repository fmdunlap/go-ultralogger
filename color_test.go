package ultralogger

import "testing"

func Test_colorize(t *testing.T) {
    type args struct {
        color Color
        msg   string
    }
    tests := []struct {
        name string
        args args
        want string
    }{
        {"RedColor", args{ColorRed, "red"}, "\033[31mred\033[0m"},
        {"GreenColor", args{ColorGreen, "green"}, "\033[32mgreen\033[0m"},
        {"YellowColor", args{ColorYellow, "yellow"}, "\033[33myellow\033[0m"},
        {"BlueColor", args{ColorBlue, "blue"}, "\033[34mblue\033[0m"},
        {"MagentaColor", args{ColorMagenta, "magenta"}, "\033[35mmagenta\033[0m"},
        {"CyanColor", args{ColorCyan, "cyan"}, "\033[36mcyan\033[0m"},
        {"GrayColor", args{ColorGray, "gray"}, "\033[2;37mgray\033[0m"},
        {"WhiteColor", args{ColorWhite, "white"}, "\033[37mwhite\033[0m"},
        {"DefaultColor", args{ColorDefault, "default"}, "\033[39mdefault\033[0m"},
        {"UnknownColor", args{Color(69), "unknown"}, "\033[39munknown\033[0m"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := colorize(tt.args.color, tt.args.msg); got != tt.want {
                t.Errorf("colorize() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_validColor(t *testing.T) {
    tests := []struct {
        name  string
        color Color
        want  bool
    }{
        {"ValidRed", ColorRed, true},
        {"ValidGreen", ColorGreen, true},
        {"ValidYellow", ColorYellow, true},
        {"ValidBlue", ColorBlue, true},
        {"ValidMagenta", ColorMagenta, true},
        {"ValidCyan", ColorCyan, true},
        {"ValidGray", ColorGray, true},
        {"ValidWhite", ColorWhite, true},
        {"ValidDefault", ColorDefault, true},
        {"invalidColor", Color(69), false},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := validColor(tt.color); got != tt.want {
                t.Errorf("validColor() = %v, want %v", got, tt.want)
            }
        })
    }
}

func Test_String(t *testing.T) {
    tests := []struct {
        name  string
        color Color
        want  string
    }{
        {"RedColor", ColorRed, "RED"},
        {"GreenColor", ColorGreen, "GREEN"},
        {"YellowColor", ColorYellow, "YELLOW"},
        {"BlueColor", ColorBlue, "BLUE"},
        {"MagentaColor", ColorMagenta, "MAGENTA"},
        {"CyanColor", ColorCyan, "CYAN"},
        {"GrayColor", ColorGray, "GRAY"},
        {"WhiteColor", ColorWhite, "WHITE"},
        {"DefaultColor", ColorDefault, "DEFAULT"},
        {"UnknownColor", Color(69), "UNKNOWN"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.color.String(); got != tt.want {
                t.Errorf("Color.String() = %v, want %v", got, tt.want)
            }
        })
    }
}
