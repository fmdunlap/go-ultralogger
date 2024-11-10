package field

import (
    "github.com/fmdunlap/go-ultralogger/level"
    "testing"
    "time"
)

type mockClock struct{}

func (c mockClock) Now() time.Time {
    return time.Date(2024, time.November, 7, 19, 30, 0, 0, time.UTC)
}

func TestDateTimeField_FieldPrinter(t *testing.T) {
    tests := []struct {
        name          string
        dateTimeField Field
        printArgs     PrintArgs
        want          string
        wantErr       bool
    }{
        {
            name: "Default",
            dateTimeField: &DateTimeField{
                dateTimeFormat: "2006-01-02 15:04:05",
                clock:          mockClock{},
            },
            printArgs: PrintArgs{
                Level: level.Info,
            },
            want: "2024-11-07 19:30:00",
        },
        {
            name: "Only Time",
            dateTimeField: &DateTimeField{
                dateTimeFormat: "15:04:05",
                clock:          mockClock{},
            },
            printArgs: PrintArgs{
                Level: level.Info,
            },
            want: "19:30:00",
        },
        {
            name: "Only Date",
            dateTimeField: &DateTimeField{
                dateTimeFormat: "2006-01-02",
                clock:          mockClock{},
            },
            printArgs: PrintArgs{
                Level: level.Info,
            },
            want: "2024-11-07",
        },
        {
            name: "Set DateTimeFormat",
            dateTimeField: (&DateTimeField{
                dateTimeFormat: "2006-01-02 15:04:05",
                clock:          mockClock{},
            }).SetDateTimeFormat("06/01/02"),
            printArgs: PrintArgs{
                Level: level.Info,
            },
            want: "24/11/07",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := tt.dateTimeField.FieldPrinter()
            if (err != nil) != tt.wantErr {
                t.Errorf("FieldPrinter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got(tt.printArgs) != tt.want {
                t.Errorf("FieldPrinter() got = %v, want %v", got(tt.printArgs), tt.want)
            }
        })
    }
}
