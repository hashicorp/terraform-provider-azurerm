package restorepoints

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorePointProperties struct {
	ConsistencyMode    *ConsistencyModeTypes       `json:"consistencyMode,omitempty"`
	ExcludeDisks       *[]ApiEntityReference       `json:"excludeDisks,omitempty"`
	InstanceView       *RestorePointInstanceView   `json:"instanceView,omitempty"`
	ProvisioningState  *string                     `json:"provisioningState,omitempty"`
	SourceMetadata     *RestorePointSourceMetadata `json:"sourceMetadata,omitempty"`
	SourceRestorePoint *ApiEntityReference         `json:"sourceRestorePoint,omitempty"`
	TimeCreated        *string                     `json:"timeCreated,omitempty"`
}

func (o *RestorePointProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *RestorePointProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}
