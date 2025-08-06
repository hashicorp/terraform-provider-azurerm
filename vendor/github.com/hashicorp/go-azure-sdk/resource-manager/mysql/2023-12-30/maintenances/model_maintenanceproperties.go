package maintenances

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceProperties struct {
	MaintenanceAvailableScheduleMaxTime *string                       `json:"maintenanceAvailableScheduleMaxTime,omitempty"`
	MaintenanceAvailableScheduleMinTime *string                       `json:"maintenanceAvailableScheduleMinTime,omitempty"`
	MaintenanceDescription              *string                       `json:"maintenanceDescription,omitempty"`
	MaintenanceEndTime                  *string                       `json:"maintenanceEndTime,omitempty"`
	MaintenanceExecutionEndTime         *string                       `json:"maintenanceExecutionEndTime,omitempty"`
	MaintenanceExecutionStartTime       *string                       `json:"maintenanceExecutionStartTime,omitempty"`
	MaintenanceStartTime                *string                       `json:"maintenanceStartTime,omitempty"`
	MaintenanceState                    *MaintenanceState             `json:"maintenanceState,omitempty"`
	MaintenanceTitle                    *string                       `json:"maintenanceTitle,omitempty"`
	MaintenanceType                     *MaintenanceType              `json:"maintenanceType,omitempty"`
	ProvisioningState                   *MaintenanceProvisioningState `json:"provisioningState,omitempty"`
}

func (o *MaintenanceProperties) GetMaintenanceAvailableScheduleMaxTimeAsTime() (*time.Time, error) {
	if o.MaintenanceAvailableScheduleMaxTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MaintenanceAvailableScheduleMaxTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MaintenanceProperties) SetMaintenanceAvailableScheduleMaxTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MaintenanceAvailableScheduleMaxTime = &formatted
}

func (o *MaintenanceProperties) GetMaintenanceAvailableScheduleMinTimeAsTime() (*time.Time, error) {
	if o.MaintenanceAvailableScheduleMinTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MaintenanceAvailableScheduleMinTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MaintenanceProperties) SetMaintenanceAvailableScheduleMinTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MaintenanceAvailableScheduleMinTime = &formatted
}

func (o *MaintenanceProperties) GetMaintenanceEndTimeAsTime() (*time.Time, error) {
	if o.MaintenanceEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MaintenanceEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MaintenanceProperties) SetMaintenanceEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MaintenanceEndTime = &formatted
}

func (o *MaintenanceProperties) GetMaintenanceExecutionEndTimeAsTime() (*time.Time, error) {
	if o.MaintenanceExecutionEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MaintenanceExecutionEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MaintenanceProperties) SetMaintenanceExecutionEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MaintenanceExecutionEndTime = &formatted
}

func (o *MaintenanceProperties) GetMaintenanceExecutionStartTimeAsTime() (*time.Time, error) {
	if o.MaintenanceExecutionStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MaintenanceExecutionStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MaintenanceProperties) SetMaintenanceExecutionStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MaintenanceExecutionStartTime = &formatted
}

func (o *MaintenanceProperties) GetMaintenanceStartTimeAsTime() (*time.Time, error) {
	if o.MaintenanceStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MaintenanceStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MaintenanceProperties) SetMaintenanceStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MaintenanceStartTime = &formatted
}
