package schedule

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/formatting"
)

type ScheduleProperties struct {
	Notes             *string            `json:"notes,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	RecurrencePattern *RecurrencePattern `json:"recurrencePattern,omitempty"`
	StartAt           *string            `json:"startAt,omitempty"`
	StopAt            string             `json:"stopAt"`
	TimeZoneId        string             `json:"timeZoneId"`
}

func (o ScheduleProperties) GetStartAtAsTime() (*time.Time, error) {
	return formatting.ParseAsDateFormat(o.StartAt, "2006-01-02T15:04:05Z07:00")
}

func (o ScheduleProperties) SetStartAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartAt = &formatted
}

func (o ScheduleProperties) GetStopAtAsTime() (*time.Time, error) {
	return formatting.ParseAsDateFormat(o.StopAt, "2006-01-02T15:04:05Z07:00")
}

func (o ScheduleProperties) SetStopAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StopAt = &formatted
}
