package synchronizationsetting

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledSynchronizationSettingProperties struct {
	CreatedAt           *string            `json:"createdAt,omitempty"`
	ProvisioningState   *ProvisioningState `json:"provisioningState,omitempty"`
	RecurrenceInterval  RecurrenceInterval `json:"recurrenceInterval"`
	SynchronizationTime string             `json:"synchronizationTime"`
	UserName            *string            `json:"userName,omitempty"`
}

func (o *ScheduledSynchronizationSettingProperties) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduledSynchronizationSettingProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o *ScheduledSynchronizationSettingProperties) GetSynchronizationTimeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.SynchronizationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduledSynchronizationSettingProperties) SetSynchronizationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SynchronizationTime = formatted
}
