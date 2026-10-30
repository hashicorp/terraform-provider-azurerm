package cloudvmclusters

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataCollectionOptions struct {
	IsDiagnosticsEventsEnabled *bool `json:"isDiagnosticsEventsEnabled,omitempty"`
	IsHealthMonitoringEnabled  *bool `json:"isHealthMonitoringEnabled,omitempty"`
	IsIncidentLogsEnabled      *bool `json:"isIncidentLogsEnabled,omitempty"`
}
