package serverendpointresource

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CloudTieringVolumeFreeSpacePolicyStatus struct {
	CurrentVolumeFreeSpacePercent  *int64  `json:"currentVolumeFreeSpacePercent,omitempty"`
	EffectiveVolumeFreeSpacePolicy *int64  `json:"effectiveVolumeFreeSpacePolicy,omitempty"`
	LastUpdatedTimestamp           *string `json:"lastUpdatedTimestamp,omitempty"`
}

func (o *CloudTieringVolumeFreeSpacePolicyStatus) GetLastUpdatedTimestampAsTime() (*time.Time, error) {
	if o.LastUpdatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *CloudTieringVolumeFreeSpacePolicyStatus) SetLastUpdatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTimestamp = &formatted
}
