package networkvirtualappliances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkVirtualAppliancePropertiesFormatNetworkProfile struct {
	NetworkInterfaceConfigurations *[]VirtualApplianceNetworkInterfaceConfiguration `json:"networkInterfaceConfigurations,omitempty"`
}
