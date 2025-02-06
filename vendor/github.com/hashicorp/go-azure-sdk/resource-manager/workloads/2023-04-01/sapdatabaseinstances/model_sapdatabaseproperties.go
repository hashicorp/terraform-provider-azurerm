package sapdatabaseinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPDatabaseProperties struct {
	DatabaseSid         *string                              `json:"databaseSid,omitempty"`
	DatabaseType        *string                              `json:"databaseType,omitempty"`
	Errors              *SAPVirtualInstanceError             `json:"errors,omitempty"`
	IPAddress           *string                              `json:"ipAddress,omitempty"`
	LoadBalancerDetails *LoadBalancerDetails                 `json:"loadBalancerDetails,omitempty"`
	ProvisioningState   *SapVirtualInstanceProvisioningState `json:"provisioningState,omitempty"`
	Status              *SAPVirtualInstanceStatus            `json:"status,omitempty"`
	Subnet              *string                              `json:"subnet,omitempty"`
	VMDetails           *[]DatabaseVMDetails                 `json:"vmDetails,omitempty"`
}
