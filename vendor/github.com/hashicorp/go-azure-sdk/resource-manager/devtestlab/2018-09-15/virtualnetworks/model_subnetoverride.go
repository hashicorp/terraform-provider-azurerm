package virtualnetworks

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SubnetOverride struct {
	LabSubnetName                      *string                                   `json:"labSubnetName,omitempty"`
	ResourceId                         *string                                   `json:"resourceId,omitempty"`
	SharedPublicIPAddressConfiguration *SubnetSharedPublicIPAddressConfiguration `json:"sharedPublicIpAddressConfiguration,omitempty"`
	UseInVMCreationPermission          *UsagePermissionType                      `json:"useInVmCreationPermission,omitempty"`
	UsePublicIPAddressPermission       *UsagePermissionType                      `json:"usePublicIpAddressPermission,omitempty"`
	VirtualNetworkPoolName             *string                                   `json:"virtualNetworkPoolName,omitempty"`
}
