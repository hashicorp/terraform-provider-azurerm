package privatednszonegroups

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateDnsZoneGroupPropertiesFormat struct {
	PrivateDnsZoneConfigs *[]PrivateDnsZoneConfig `json:"privateDnsZoneConfigs,omitempty"`
	ProvisioningState     *ProvisioningState      `json:"provisioningState,omitempty"`
}
