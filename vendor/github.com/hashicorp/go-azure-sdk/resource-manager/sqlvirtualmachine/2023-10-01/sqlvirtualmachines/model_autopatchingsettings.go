package sqlvirtualmachines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoPatchingSettings struct {
	AdditionalVMPatch             *AdditionalVMPatch `json:"additionalVmPatch,omitempty"`
	DayOfWeek                     *DayOfWeek         `json:"dayOfWeek,omitempty"`
	Enable                        *bool              `json:"enable,omitempty"`
	MaintenanceWindowDuration     *int64             `json:"maintenanceWindowDuration,omitempty"`
	MaintenanceWindowStartingHour *int64             `json:"maintenanceWindowStartingHour,omitempty"`
}
