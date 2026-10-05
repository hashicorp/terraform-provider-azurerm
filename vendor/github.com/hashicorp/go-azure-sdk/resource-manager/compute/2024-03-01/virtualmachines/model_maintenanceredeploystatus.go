package virtualmachines

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MaintenanceRedeployStatus struct {
	IsCustomerInitiatedMaintenanceAllowed *bool                                `json:"isCustomerInitiatedMaintenanceAllowed,omitempty"`
	LastOperationMessage                  *string                              `json:"lastOperationMessage,omitempty"`
	LastOperationResultCode               *MaintenanceOperationResultCodeTypes `json:"lastOperationResultCode,omitempty"`
	MaintenanceWindowEndTime              *string                              `json:"maintenanceWindowEndTime,omitempty"`
	MaintenanceWindowStartTime            *string                              `json:"maintenanceWindowStartTime,omitempty"`
	PreMaintenanceWindowEndTime           *string                              `json:"preMaintenanceWindowEndTime,omitempty"`
	PreMaintenanceWindowStartTime         *string                              `json:"preMaintenanceWindowStartTime,omitempty"`
}

func (o *MaintenanceRedeployStatus) GetMaintenanceWindowEndTimeAsTime() (*time.Time, error) {
	if o.MaintenanceWindowEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MaintenanceWindowEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MaintenanceRedeployStatus) SetMaintenanceWindowEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MaintenanceWindowEndTime = &formatted
}

func (o *MaintenanceRedeployStatus) GetMaintenanceWindowStartTimeAsTime() (*time.Time, error) {
	if o.MaintenanceWindowStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.MaintenanceWindowStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MaintenanceRedeployStatus) SetMaintenanceWindowStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.MaintenanceWindowStartTime = &formatted
}

func (o *MaintenanceRedeployStatus) GetPreMaintenanceWindowEndTimeAsTime() (*time.Time, error) {
	if o.PreMaintenanceWindowEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PreMaintenanceWindowEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MaintenanceRedeployStatus) SetPreMaintenanceWindowEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PreMaintenanceWindowEndTime = &formatted
}

func (o *MaintenanceRedeployStatus) GetPreMaintenanceWindowStartTimeAsTime() (*time.Time, error) {
	if o.PreMaintenanceWindowStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PreMaintenanceWindowStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *MaintenanceRedeployStatus) SetPreMaintenanceWindowStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PreMaintenanceWindowStartTime = &formatted
}
