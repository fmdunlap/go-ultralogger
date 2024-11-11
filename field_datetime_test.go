package ultralogger

import (
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
        llCtx         LogLineContext
        want          string
        wantErr       bool
    }{
        {
            name: "Default",
            dateTimeField: &FieldDateTime{
                dateTimeFormat: "2006-01-02 15:04:05",
                clock:          mockClock{},
            },
            llCtx: LogLineContext{
                Level: Info,
            },
            want: "2024-11-07 19:30:00",
        },
        {
            name: "Only Time",
            dateTimeField: &FieldDateTime{
                dateTimeFormat: "15:04:05",
                clock:          mockClock{},
            },
            llCtx: LogLineContext{
                Level: Info,
            },
            want: "19:30:00",
        },
        {
            name: "Only Date",
            dateTimeField: &FieldDateTime{
                dateTimeFormat: "2006-01-02",
                clock:          mockClock{},
            },
            llCtx: LogLineContext{
                Level: Info,
            },
            want: "2024-11-07",
        },
        {
            name: "Set DateTimeFormat",
            dateTimeField: (&FieldDateTime{
                dateTimeFormat: "2006-01-02 15:04:05",
                clock:          mockClock{},
            }).SetDateTimeFormat("06/01/02"),
            llCtx: LogLineContext{
                Level: Info,
            },
            want: "24/11/07",
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := tt.dateTimeField.FieldFormatter()
            if (err != nil) != tt.wantErr {
                t.Errorf("FieldFormatter() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got(tt.llCtx, OutputFormatText, nil).Data != tt.want {
                t.Errorf("FieldFormatter() got = %v, want %v", got(tt.llCtx, OutputFormatText, nil), tt.want)
            }
        })
    }
}
