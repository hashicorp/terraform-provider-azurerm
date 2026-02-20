package cosmosdb

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PercentileMetricValue struct {
	Average   *float64 `json:"average,omitempty"`
	Count     *float64 `json:"_count,omitempty"`
	Maximum   *float64 `json:"maximum,omitempty"`
	Minimum   *float64 `json:"minimum,omitempty"`
	P10       *float64 `json:"P10,omitempty"`
	P25       *float64 `json:"P25,omitempty"`
	P50       *float64 `json:"P50,omitempty"`
	P75       *float64 `json:"P75,omitempty"`
	P90       *float64 `json:"P90,omitempty"`
	P95       *float64 `json:"P95,omitempty"`
	P99       *float64 `json:"P99,omitempty"`
	Timestamp *string  `json:"timestamp,omitempty"`
	Total     *float64 `json:"total,omitempty"`
}

func (o *PercentileMetricValue) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *PercentileMetricValue) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}
