package serverendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudTieringSpaceSavings struct {
	CachedSizeBytes      *int64  `json:"cachedSizeBytes,omitempty"`
	LastUpdatedTimestamp *string `json:"lastUpdatedTimestamp,omitempty"`
	SpaceSavingsBytes    *int64  `json:"spaceSavingsBytes,omitempty"`
	SpaceSavingsPercent  *int64  `json:"spaceSavingsPercent,omitempty"`
	TotalSizeCloudBytes  *int64  `json:"totalSizeCloudBytes,omitempty"`
	VolumeSizeBytes      *int64  `json:"volumeSizeBytes,omitempty"`
}

func (o *CloudTieringSpaceSavings) GetLastUpdatedTimestampAsTime() (*time.Time, error) {
	if o.LastUpdatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudTieringSpaceSavings) SetLastUpdatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTimestamp = &formatted
}
