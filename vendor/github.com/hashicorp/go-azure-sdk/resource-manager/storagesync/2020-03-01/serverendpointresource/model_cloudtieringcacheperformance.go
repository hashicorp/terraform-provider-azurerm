package serverendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudTieringCachePerformance struct {
	CacheHitBytes        *int64  `json:"cacheHitBytes,omitempty"`
	CacheHitBytesPercent *int64  `json:"cacheHitBytesPercent,omitempty"`
	CacheMissBytes       *int64  `json:"cacheMissBytes,omitempty"`
	LastUpdatedTimestamp *string `json:"lastUpdatedTimestamp,omitempty"`
}

func (o *CloudTieringCachePerformance) GetLastUpdatedTimestampAsTime() (*time.Time, error) {
	if o.LastUpdatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudTieringCachePerformance) SetLastUpdatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTimestamp = &formatted
}
