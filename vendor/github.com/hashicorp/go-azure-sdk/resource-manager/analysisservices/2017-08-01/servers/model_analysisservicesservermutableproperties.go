package servers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AnalysisServicesServerMutableProperties struct {
	AsAdministrators        *ServerAdministrators `json:"asAdministrators,omitempty"`
	BackupBlobContainerUri  *string               `json:"backupBlobContainerUri,omitempty"`
	GatewayDetails          *GatewayDetails       `json:"gatewayDetails,omitempty"`
	IpV4FirewallSettings    *IPv4FirewallSettings `json:"ipV4FirewallSettings,omitempty"`
	ManagedMode             *ManagedMode          `json:"managedMode,omitempty"`
	QuerypoolConnectionMode *ConnectionMode       `json:"querypoolConnectionMode,omitempty"`
	ServerMonitorMode       *ServerMonitorMode    `json:"serverMonitorMode,omitempty"`
}
