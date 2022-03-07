package heatmaps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type HeatMapProperties struct {
	EndTime      *string            `json:"endTime,omitempty"`
	Endpoints    *[]HeatMapEndpoint `json:"endpoints,omitempty"`
	StartTime    *string            `json:"startTime,omitempty"`
	TrafficFlows *[]TrafficFlow     `json:"trafficFlows,omitempty"`
}

func (o HeatMapProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o HeatMapProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o HeatMapProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o HeatMapProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
