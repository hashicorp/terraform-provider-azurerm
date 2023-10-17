package cloudendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudEndpointChangeEnumerationActivity struct {
	DeletesProgressPercent    *int64                                          `json:"deletesProgressPercent,omitempty"`
	LastUpdatedTimestamp      *string                                         `json:"lastUpdatedTimestamp,omitempty"`
	MinutesRemaining          *int64                                          `json:"minutesRemaining,omitempty"`
	OperationState            *CloudEndpointChangeEnumerationActivityState    `json:"operationState,omitempty"`
	ProcessedDirectoriesCount *int64                                          `json:"processedDirectoriesCount,omitempty"`
	ProcessedFilesCount       *int64                                          `json:"processedFilesCount,omitempty"`
	ProgressPercent           *int64                                          `json:"progressPercent,omitempty"`
	StartedTimestamp          *string                                         `json:"startedTimestamp,omitempty"`
	StatusCode                *int64                                          `json:"statusCode,omitempty"`
	TotalCountsState          *CloudEndpointChangeEnumerationTotalCountsState `json:"totalCountsState,omitempty"`
	TotalDirectoriesCount     *int64                                          `json:"totalDirectoriesCount,omitempty"`
	TotalFilesCount           *int64                                          `json:"totalFilesCount,omitempty"`
	TotalSizeBytes            *int64                                          `json:"totalSizeBytes,omitempty"`
}

func (o *CloudEndpointChangeEnumerationActivity) GetLastUpdatedTimestampAsTime() (*time.Time, error) {
	if o.LastUpdatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudEndpointChangeEnumerationActivity) SetLastUpdatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTimestamp = &formatted
}

func (o *CloudEndpointChangeEnumerationActivity) GetStartedTimestampAsTime() (*time.Time, error) {
	if o.StartedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudEndpointChangeEnumerationActivity) SetStartedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartedTimestamp = &formatted
}
