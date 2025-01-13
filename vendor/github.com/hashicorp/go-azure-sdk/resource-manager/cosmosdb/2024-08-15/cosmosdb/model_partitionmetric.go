package cosmosdb

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PartitionMetric struct {
	EndTime             *string        `json:"endTime,omitempty"`
	MetricValues        *[]MetricValue `json:"metricValues,omitempty"`
	Name                *MetricName    `json:"name,omitempty"`
	PartitionId         *string        `json:"partitionId,omitempty"`
	PartitionKeyRangeId *string        `json:"partitionKeyRangeId,omitempty"`
	StartTime           *string        `json:"startTime,omitempty"`
	TimeGrain           *string        `json:"timeGrain,omitempty"`
	Unit                *UnitType      `json:"unit,omitempty"`
}

func (o *PartitionMetric) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *PartitionMetric) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *PartitionMetric) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *PartitionMetric) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
