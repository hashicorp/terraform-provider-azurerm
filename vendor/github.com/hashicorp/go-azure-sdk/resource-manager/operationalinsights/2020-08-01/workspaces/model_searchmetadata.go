package workspaces

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SearchMetadata struct {
	AggregatedGroupingFields *string               `json:"aggregatedGroupingFields,omitempty"`
	AggregatedValueField     *string               `json:"aggregatedValueField,omitempty"`
	CoreSummaries            *[]CoreSummary        `json:"coreSummaries,omitempty"`
	ETag                     *string               `json:"eTag,omitempty"`
	Id                       *string               `json:"id,omitempty"`
	LastUpdated              *string               `json:"lastUpdated,omitempty"`
	Max                      *int64                `json:"max,omitempty"`
	RequestId                *string               `json:"requestId,omitempty"`
	RequestTime              *int64                `json:"requestTime,omitempty"`
	ResultType               *string               `json:"resultType,omitempty"`
	Schema                   *SearchMetadataSchema `json:"schema,omitempty"`
	Sort                     *[]SearchSort         `json:"sort,omitempty"`
	StartTime                *string               `json:"startTime,omitempty"`
	Status                   *string               `json:"status,omitempty"`
	Sum                      *int64                `json:"sum,omitempty"`
	Top                      *int64                `json:"top,omitempty"`
	Total                    *int64                `json:"total,omitempty"`
}

func (o *SearchMetadata) GetLastUpdatedAsTime() (*time.Time, error) {
	if o.LastUpdated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdated, "2006-01-02T15:04:05Z07:00")
}

func (o *SearchMetadata) SetLastUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdated = &formatted
}

func (o *SearchMetadata) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SearchMetadata) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
