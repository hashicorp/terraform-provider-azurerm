package maintenances

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenancePropertiesForUpdate struct {
	MaintenanceStartTime *string `json:"maintenanceStartTime,omitempty"`
}

func (o *MaintenancePropertiesForUpdate) GetMaintenanceStartTimeAsTime() (*time.Time, error) {
	if o.MaintenanceStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MaintenanceStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MaintenancePropertiesForUpdate) SetMaintenanceStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MaintenanceStartTime = &formatted
}
