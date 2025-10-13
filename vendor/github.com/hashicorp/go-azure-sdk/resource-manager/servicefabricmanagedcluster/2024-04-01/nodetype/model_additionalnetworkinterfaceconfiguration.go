package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AdditionalNetworkInterfaceConfiguration struct {
	DscpConfiguration           *SubResource      `json:"dscpConfiguration,omitempty"`
	EnableAcceleratedNetworking *bool             `json:"enableAcceleratedNetworking,omitempty"`
	IPConfigurations            []IPConfiguration `json:"ipConfigurations"`
	Name                        string            `json:"name"`
}
