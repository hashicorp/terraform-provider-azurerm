package fleets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ComputeProfile struct {
	AdditionalVirtualMachineCapabilities *AdditionalCapabilities   `json:"additionalVirtualMachineCapabilities,omitempty"`
	BaseVirtualMachineProfile            BaseVirtualMachineProfile `json:"baseVirtualMachineProfile"`
	ComputeApiVersion                    *string                   `json:"computeApiVersion,omitempty"`
	PlatformFaultDomainCount             *int64                    `json:"platformFaultDomainCount,omitempty"`
}
