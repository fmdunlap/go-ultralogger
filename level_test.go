package ultralogger

import (
    "reflect"
    "testing"
)

func TestAllLevels(t *testing.T) {
    tests := []struct {
        name string
        want []Level
    }{
        {
            "AllLevels",
            []Level{
                DebugLevel,
                InfoLevel,
                WarnLevel,
                ErrorLevel,
                PanicLevel,
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := AllLevels(); !reflect.DeepEqual(got, tt.want) {
                t.Errorf("AllLevels() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestLevel_String(t *testing.T) {
    tests := []struct {
        name string
        l    Level
        want string
    }{
        {"DebugLevel", DebugLevel, "DEBUG"},
        {"InfoLevel", InfoLevel, "INFO"},
        {"WarnLevel", WarnLevel, "WARN"},
        {"ErrorLevel", ErrorLevel, "ERROR"},
        {"PanicLevel", PanicLevel, "PANIC"},
        {"UnknownLevel", Level(42), "UNKNOWN"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.l.String(); got != tt.want {
                t.Errorf("String() = %v, want %v", got, tt.want)
            }
        })
    }
}

func TestParseLevel(t *testing.T) {
    type args struct {
        levelStr string
    }
    tests := []struct {
        name    string
        args    args
        want    Level
        wantErr bool
    }{
        {"DebugLevel", args{"debug"}, DebugLevel, false},
        {"InfoLevel", args{"info"}, InfoLevel, false},
        {"WarnLevel", args{"warn"}, WarnLevel, false},
        {"ErrorLevel", args{"error"}, ErrorLevel, false},
        {"PanicLevel", args{"panic"}, PanicLevel, false},
        {"InvalidLevel", args{"invalid"}, 0, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseLevel(tt.args.levelStr)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseLevel() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ParseLevel() got = %v, want %v", got, tt.want)
            }
        })
    }
}
