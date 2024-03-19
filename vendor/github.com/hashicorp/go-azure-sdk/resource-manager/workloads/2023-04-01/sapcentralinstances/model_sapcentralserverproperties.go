package sapcentralinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPCentralServerProperties struct {
	EnqueueReplicationServerProperties *EnqueueReplicationServerProperties  `json:"enqueueReplicationServerProperties,omitempty"`
	EnqueueServerProperties            *EnqueueServerProperties             `json:"enqueueServerProperties,omitempty"`
	Errors                             *SAPVirtualInstanceError             `json:"errors,omitempty"`
	GatewayServerProperties            *GatewayServerProperties             `json:"gatewayServerProperties,omitempty"`
	Health                             *SAPHealthState                      `json:"health,omitempty"`
	InstanceNo                         *string                              `json:"instanceNo,omitempty"`
	KernelPatch                        *string                              `json:"kernelPatch,omitempty"`
	KernelVersion                      *string                              `json:"kernelVersion,omitempty"`
	LoadBalancerDetails                *LoadBalancerDetails                 `json:"loadBalancerDetails,omitempty"`
	MessageServerProperties            *MessageServerProperties             `json:"messageServerProperties,omitempty"`
	ProvisioningState                  *SapVirtualInstanceProvisioningState `json:"provisioningState,omitempty"`
	Status                             *SAPVirtualInstanceStatus            `json:"status,omitempty"`
	Subnet                             *string                              `json:"subnet,omitempty"`
	VMDetails                          *[]CentralServerVMDetails            `json:"vmDetails,omitempty"`
}
