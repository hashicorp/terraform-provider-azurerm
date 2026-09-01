package virtualmachinescalesetrollingupgrades

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RollingUpgradeRunningStatus struct {
	Code           *RollingUpgradeStatusCode `json:"code,omitempty"`
	LastAction     *RollingUpgradeActionType `json:"lastAction,omitempty"`
	LastActionTime *string                   `json:"lastActionTime,omitempty"`
	StartTime      *string                   `json:"startTime,omitempty"`
}

func (o *RollingUpgradeRunningStatus) GetLastActionTimeAsTime() (*time.Time, error) {
	if o.LastActionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastActionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RollingUpgradeRunningStatus) SetLastActionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastActionTime = &formatted
}

func (o *RollingUpgradeRunningStatus) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RollingUpgradeRunningStatus) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
