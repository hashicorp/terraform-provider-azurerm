package networkprofiles

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerNetworkInterfaceConfigurationPropertiesFormat struct {
	ContainerNetworkInterfaces *[]SubResource            `json:"containerNetworkInterfaces,omitempty"`
	IPConfigurations           *[]IPConfigurationProfile `json:"ipConfigurations,omitempty"`
	ProvisioningState          *ProvisioningState        `json:"provisioningState,omitempty"`
}
