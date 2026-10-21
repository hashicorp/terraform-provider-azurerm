package virtualmachinescalesetvms

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NetworkProfile struct {
	NetworkApiVersion              *NetworkApiVersion                             `json:"networkApiVersion,omitempty"`
	NetworkInterfaceConfigurations *[]VirtualMachineNetworkInterfaceConfiguration `json:"networkInterfaceConfigurations,omitempty"`
	NetworkInterfaces              *[]NetworkInterfaceReference                   `json:"networkInterfaces,omitempty"`
}
