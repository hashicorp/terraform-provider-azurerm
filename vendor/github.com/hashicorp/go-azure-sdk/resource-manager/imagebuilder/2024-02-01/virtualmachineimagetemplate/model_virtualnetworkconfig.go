package virtualmachineimagetemplate

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkConfig struct {
	ContainerInstanceSubnetId *string `json:"containerInstanceSubnetId,omitempty"`
	ProxyVMSize               *string `json:"proxyVmSize,omitempty"`
	SubnetId                  *string `json:"subnetId,omitempty"`
}
