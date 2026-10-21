package profiles

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MetricsResponse struct {
	DateTimeBegin *string                      `json:"dateTimeBegin,omitempty"`
	DateTimeEnd   *string                      `json:"dateTimeEnd,omitempty"`
	Granularity   *MetricsGranularity          `json:"granularity,omitempty"`
	Series        *[]MetricsResponseSeriesItem `json:"series,omitempty"`
}

func (o *MetricsResponse) GetDateTimeBeginAsTime() (*time.Time, error) {
	if o.DateTimeBegin == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DateTimeBegin, "2006-01-02T15:04:05Z07:00")
}

func (o *MetricsResponse) SetDateTimeBeginAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DateTimeBegin = &formatted
}

func (o *MetricsResponse) GetDateTimeEndAsTime() (*time.Time, error) {
	if o.DateTimeEnd == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DateTimeEnd, "2006-01-02T15:04:05Z07:00")
}

func (o *MetricsResponse) SetDateTimeEndAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DateTimeEnd = &formatted
}
