package workspaces

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UsageMetric struct {
	CurrentValue  *float64    `json:"currentValue,omitempty"`
	Limit         *float64    `json:"limit,omitempty"`
	Name          *MetricName `json:"name,omitempty"`
	NextResetTime *string     `json:"nextResetTime,omitempty"`
	QuotaPeriod   *string     `json:"quotaPeriod,omitempty"`
	Unit          *string     `json:"unit,omitempty"`
}

func (o *UsageMetric) GetNextResetTimeAsTime() (*time.Time, error) {
	if o.NextResetTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NextResetTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UsageMetric) SetNextResetTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NextResetTime = &formatted
}
