package cloudservicepublicipaddresses

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkSecurityGroupPropertiesFormat struct {
	DefaultSecurityRules *[]SecurityRule     `json:"defaultSecurityRules,omitempty"`
	FlowLogs             *[]FlowLog          `json:"flowLogs,omitempty"`
	FlushConnection      *bool               `json:"flushConnection,omitempty"`
	NetworkInterfaces    *[]NetworkInterface `json:"networkInterfaces,omitempty"`
	ProvisioningState    *ProvisioningState  `json:"provisioningState,omitempty"`
	ResourceGuid         *string             `json:"resourceGuid,omitempty"`
	SecurityRules        *[]SecurityRule     `json:"securityRules,omitempty"`
	Subnets              *[]Subnet           `json:"subnets,omitempty"`
}
