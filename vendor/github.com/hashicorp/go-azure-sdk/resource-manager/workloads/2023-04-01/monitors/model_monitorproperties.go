package monitors

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorProperties struct {
	AppLocation                       *string                           `json:"appLocation,omitempty"`
	Errors                            *Error                            `json:"errors,omitempty"`
	LogAnalyticsWorkspaceArmId        *string                           `json:"logAnalyticsWorkspaceArmId,omitempty"`
	ManagedResourceGroupConfiguration *ManagedRGConfiguration           `json:"managedResourceGroupConfiguration,omitempty"`
	MonitorSubnet                     *string                           `json:"monitorSubnet,omitempty"`
	MsiArmId                          *string                           `json:"msiArmId,omitempty"`
	ProvisioningState                 *WorkloadMonitorProvisioningState `json:"provisioningState,omitempty"`
	RoutingPreference                 *RoutingPreference                `json:"routingPreference,omitempty"`
	StorageAccountArmId               *string                           `json:"storageAccountArmId,omitempty"`
	ZoneRedundancyPreference          *string                           `json:"zoneRedundancyPreference,omitempty"`
}
