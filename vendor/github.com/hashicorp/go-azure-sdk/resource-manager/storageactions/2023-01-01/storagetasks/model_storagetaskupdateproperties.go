package storagetasks

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StorageTaskUpdateProperties struct {
	Action            *StorageTaskAction `json:"action,omitempty"`
	CreationTimeInUtc *string            `json:"creationTimeInUtc,omitempty"`
	Description       *string            `json:"description,omitempty"`
	Enabled           *bool              `json:"enabled,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
	TaskVersion       *int64             `json:"taskVersion,omitempty"`
}

func (o *StorageTaskUpdateProperties) GetCreationTimeInUtcAsTime() (*time.Time, error) {
	if o.CreationTimeInUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTimeInUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *StorageTaskUpdateProperties) SetCreationTimeInUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTimeInUtc = &formatted
}
