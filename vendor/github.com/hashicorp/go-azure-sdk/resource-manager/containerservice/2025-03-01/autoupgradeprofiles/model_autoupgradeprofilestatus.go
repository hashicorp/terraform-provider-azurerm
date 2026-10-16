package autoupgradeprofiles

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoUpgradeProfileStatus struct {
	LastTriggerError           *ErrorDetail                  `json:"lastTriggerError,omitempty"`
	LastTriggerStatus          *AutoUpgradeLastTriggerStatus `json:"lastTriggerStatus,omitempty"`
	LastTriggerUpgradeVersions *[]string                     `json:"lastTriggerUpgradeVersions,omitempty"`
	LastTriggeredAt            *string                       `json:"lastTriggeredAt,omitempty"`
}

func (o *AutoUpgradeProfileStatus) GetLastTriggeredAtAsTime() (*time.Time, error) {
	if o.LastTriggeredAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastTriggeredAt, "2006-01-02T15:04:05Z07:00")
}

func (o *AutoUpgradeProfileStatus) SetLastTriggeredAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastTriggeredAt = &formatted
}
