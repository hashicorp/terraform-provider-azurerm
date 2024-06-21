package virtualmachinescalesets

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualMachineScaleSetNetworkConfiguration struct {
	Name       string                                                `json:"name"`
	Properties *VirtualMachineScaleSetNetworkConfigurationProperties `json:"properties,omitempty"`
}
