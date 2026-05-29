package updateruns

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateRunProperties struct {
	Duration          *string                   `json:"duration,omitempty"`
	LastUpdatedTime   *string                   `json:"lastUpdatedTime,omitempty"`
	Progress          *Step                     `json:"progress,omitempty"`
	ProvisioningState *ProvisioningState        `json:"provisioningState,omitempty"`
	State             *UpdateRunPropertiesState `json:"state,omitempty"`
	TimeStarted       *string                   `json:"timeStarted,omitempty"`
}

func (o *UpdateRunProperties) GetLastUpdatedTimeAsTime() (*time.Time, error) {
	if o.LastUpdatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateRunProperties) SetLastUpdatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTime = &formatted
}

func (o *UpdateRunProperties) GetTimeStartedAsTime() (*time.Time, error) {
	if o.TimeStarted == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeStarted, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateRunProperties) SetTimeStartedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeStarted = &formatted
}
