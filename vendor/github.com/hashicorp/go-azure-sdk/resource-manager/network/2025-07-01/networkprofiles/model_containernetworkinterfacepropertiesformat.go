package networkprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerNetworkInterfacePropertiesFormat struct {
	Container                              *SubResource                                `json:"container,omitempty"`
	ContainerNetworkInterfaceConfiguration *ContainerNetworkInterfaceConfiguration     `json:"containerNetworkInterfaceConfiguration,omitempty"`
	IPConfigurations                       *[]ContainerNetworkInterfaceIPConfiguration `json:"ipConfigurations,omitempty"`
	ProvisioningState                      *ProvisioningState                          `json:"provisioningState,omitempty"`
}
