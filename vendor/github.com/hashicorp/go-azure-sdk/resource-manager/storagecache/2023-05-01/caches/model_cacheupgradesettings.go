package caches

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CacheUpgradeSettings struct {
	ScheduledTime          *string `json:"scheduledTime,omitempty"`
	UpgradeScheduleEnabled *bool   `json:"upgradeScheduleEnabled,omitempty"`
}

func (o *CacheUpgradeSettings) GetScheduledTimeAsTime() (*time.Time, error) {
	if o.ScheduledTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ScheduledTime, "2006-01-02T15:04:05Z07:00")
}

func (o *CacheUpgradeSettings) SetScheduledTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ScheduledTime = &formatted
}
