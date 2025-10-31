package listdeletedrunbooks

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeletedRunbookProperties struct {
	CreationTime       *string `json:"creationTime,omitempty"`
	DeletionTime       *string `json:"deletionTime,omitempty"`
	RunbookId          *string `json:"runbookId,omitempty"`
	RunbookType        *string `json:"runbookType,omitempty"`
	Runtime            *string `json:"runtime,omitempty"`
	RuntimeEnvironment *string `json:"runtimeEnvironment,omitempty"`
}

func (o *DeletedRunbookProperties) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DeletedRunbookProperties) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *DeletedRunbookProperties) GetDeletionTimeAsTime() (*time.Time, error) {
	if o.DeletionTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DeletionTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DeletedRunbookProperties) SetDeletionTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DeletionTime = &formatted
}
