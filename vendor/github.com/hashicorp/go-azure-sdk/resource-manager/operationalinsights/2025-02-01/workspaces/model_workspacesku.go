package workspaces

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceSku struct {
	CapacityReservationLevel *CapacityReservationLevel `json:"capacityReservationLevel,omitempty"`
	LastSkuUpdate            *string                   `json:"lastSkuUpdate,omitempty"`
	Name                     WorkspaceSkuNameEnum      `json:"name"`
}

func (o *WorkspaceSku) GetLastSkuUpdateAsTime() (*time.Time, error) {
	if o.LastSkuUpdate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSkuUpdate, "2006-01-02T15:04:05Z07:00")
}

func (o *WorkspaceSku) SetLastSkuUpdateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSkuUpdate = &formatted
}
