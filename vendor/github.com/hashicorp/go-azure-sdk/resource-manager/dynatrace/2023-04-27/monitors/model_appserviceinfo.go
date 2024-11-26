package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppServiceInfo struct {
	AutoUpdateSetting *AutoUpdateSetting `json:"autoUpdateSetting,omitempty"`
	AvailabilityState *AvailabilityState `json:"availabilityState,omitempty"`
	HostGroup         *string            `json:"hostGroup,omitempty"`
	HostName          *string            `json:"hostName,omitempty"`
	LogModule         *LogModule         `json:"logModule,omitempty"`
	MonitoringType    *MonitoringType    `json:"monitoringType,omitempty"`
	ResourceId        *string            `json:"resourceId,omitempty"`
	UpdateStatus      *UpdateStatus      `json:"updateStatus,omitempty"`
	Version           *string            `json:"version,omitempty"`
}
