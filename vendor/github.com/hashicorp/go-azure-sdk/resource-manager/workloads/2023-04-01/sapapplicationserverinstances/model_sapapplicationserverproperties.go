package sapapplicationserverinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SAPApplicationServerProperties struct {
	Errors              *SAPVirtualInstanceError             `json:"errors,omitempty"`
	GatewayPort         *int64                               `json:"gatewayPort,omitempty"`
	Health              *SAPHealthState                      `json:"health,omitempty"`
	Hostname            *string                              `json:"hostname,omitempty"`
	IPAddress           *string                              `json:"ipAddress,omitempty"`
	IcmHTTPPort         *int64                               `json:"icmHttpPort,omitempty"`
	IcmHTTPSPort        *int64                               `json:"icmHttpsPort,omitempty"`
	InstanceNo          *string                              `json:"instanceNo,omitempty"`
	KernelPatch         *string                              `json:"kernelPatch,omitempty"`
	KernelVersion       *string                              `json:"kernelVersion,omitempty"`
	LoadBalancerDetails *LoadBalancerDetails                 `json:"loadBalancerDetails,omitempty"`
	ProvisioningState   *SapVirtualInstanceProvisioningState `json:"provisioningState,omitempty"`
	Status              *SAPVirtualInstanceStatus            `json:"status,omitempty"`
	Subnet              *string                              `json:"subnet,omitempty"`
	VMDetails           *[]ApplicationServerVMDetails        `json:"vmDetails,omitempty"`
}
