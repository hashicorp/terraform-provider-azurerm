package attacheddatanetwork

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AttachedDataNetworkPropertiesFormat struct {
	DnsAddresses                         []string            `json:"dnsAddresses"`
	NaptConfiguration                    *NaptConfiguration  `json:"naptConfiguration,omitempty"`
	ProvisioningState                    *ProvisioningState  `json:"provisioningState,omitempty"`
	UserEquipmentAddressPoolPrefix       *[]string           `json:"userEquipmentAddressPoolPrefix,omitempty"`
	UserEquipmentStaticAddressPoolPrefix *[]string           `json:"userEquipmentStaticAddressPoolPrefix,omitempty"`
	UserPlaneDataInterface               InterfaceProperties `json:"userPlaneDataInterface"`
}
