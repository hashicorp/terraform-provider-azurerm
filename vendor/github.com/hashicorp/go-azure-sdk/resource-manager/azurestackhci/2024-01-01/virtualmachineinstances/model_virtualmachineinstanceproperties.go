package virtualmachineinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineInstanceProperties struct {
	GuestAgentInstallStatus *GuestAgentInstallStatus                         `json:"guestAgentInstallStatus,omitempty"`
	HTTPProxyConfig         *HTTPProxyConfiguration                          `json:"httpProxyConfig,omitempty"`
	HardwareProfile         *VirtualMachineInstancePropertiesHardwareProfile `json:"hardwareProfile,omitempty"`
	InstanceView            *VirtualMachineInstanceView                      `json:"instanceView,omitempty"`
	NetworkProfile          *VirtualMachineInstancePropertiesNetworkProfile  `json:"networkProfile,omitempty"`
	OsProfile               *VirtualMachineInstancePropertiesOsProfile       `json:"osProfile,omitempty"`
	ProvisioningState       *ProvisioningStateEnum                           `json:"provisioningState,omitempty"`
	ResourceUid             *string                                          `json:"resourceUid,omitempty"`
	SecurityProfile         *VirtualMachineInstancePropertiesSecurityProfile `json:"securityProfile,omitempty"`
	Status                  *VirtualMachineInstanceStatus                    `json:"status,omitempty"`
	StorageProfile          *VirtualMachineInstancePropertiesStorageProfile  `json:"storageProfile,omitempty"`
	VMId                    *string                                          `json:"vmId,omitempty"`
}
