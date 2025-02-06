package machines

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMwareNetworkAdapter struct {
	IPAddressList *[]string `json:"ipAddressList,omitempty"`
	IPAddressType *string   `json:"ipAddressType,omitempty"`
	Label         *string   `json:"label,omitempty"`
	MacAddress    *string   `json:"macAddress,omitempty"`
	NetworkName   *string   `json:"networkName,omitempty"`
	NicId         *string   `json:"nicId,omitempty"`
}
