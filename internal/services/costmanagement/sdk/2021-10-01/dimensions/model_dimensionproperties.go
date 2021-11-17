package dimensions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

type DimensionProperties struct {
	Category        *string   `json:"category,omitempty"`
	Data            *[]string `json:"data,omitempty"`
	Description     *string   `json:"description,omitempty"`
	FilterEnabled   *bool     `json:"filterEnabled,omitempty"`
	GroupingEnabled *bool     `json:"groupingEnabled,omitempty"`
	NextLink        *string   `json:"nextLink,omitempty"`
	Total           *int64    `json:"total,omitempty"`
	UsageEnd        *string   `json:"usageEnd,omitempty"`
	UsageStart      *string   `json:"usageStart,omitempty"`
}

func (o DimensionProperties) GetUsageEndAsTime() (*time.Time, error) {
	if o.UsageEnd == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UsageEnd, "2006-01-02T15:04:05Z07:00")
}

func (o DimensionProperties) SetUsageEndAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UsageEnd = &formatted
}

func (o DimensionProperties) GetUsageStartAsTime() (*time.Time, error) {
	if o.UsageStart == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UsageStart, "2006-01-02T15:04:05Z07:00")
}

func (o DimensionProperties) SetUsageStartAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UsageStart = &formatted
}
