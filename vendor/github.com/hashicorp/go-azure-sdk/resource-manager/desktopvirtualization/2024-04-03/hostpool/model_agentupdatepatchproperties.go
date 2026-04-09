package hostpool

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentUpdatePatchProperties struct {
	MaintenanceWindowTimeZone *string                             `json:"maintenanceWindowTimeZone,omitempty"`
	MaintenanceWindows        *[]MaintenanceWindowPatchProperties `json:"maintenanceWindows,omitempty"`
	Type                      *SessionHostComponentUpdateType     `json:"type,omitempty"`
	UseSessionHostLocalTime   *bool                               `json:"useSessionHostLocalTime,omitempty"`
}
