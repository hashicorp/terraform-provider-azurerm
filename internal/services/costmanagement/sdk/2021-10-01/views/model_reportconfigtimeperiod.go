package views

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type ReportConfigTimePeriod struct {
	From string `json:"from"`
	To   string `json:"to"`
}

func (o ReportConfigTimePeriod) GetFromAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.From, "2006-01-02T15:04:05Z07:00")
}

func (o ReportConfigTimePeriod) SetFromAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.From = formatted
}

func (o ReportConfigTimePeriod) GetToAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.To, "2006-01-02T15:04:05Z07:00")
}

func (o ReportConfigTimePeriod) SetToAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.To = formatted
}
