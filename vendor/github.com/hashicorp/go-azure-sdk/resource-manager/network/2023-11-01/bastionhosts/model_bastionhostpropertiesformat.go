package bastionhosts

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type BastionHostPropertiesFormat struct {
	DisableCopyPaste    *bool                                   `json:"disableCopyPaste,omitempty"`
	DnsName             *string                                 `json:"dnsName,omitempty"`
	EnableFileCopy      *bool                                   `json:"enableFileCopy,omitempty"`
	EnableIPConnect     *bool                                   `json:"enableIpConnect,omitempty"`
	EnableKerberos      *bool                                   `json:"enableKerberos,omitempty"`
	EnableShareableLink *bool                                   `json:"enableShareableLink,omitempty"`
	EnableTunneling     *bool                                   `json:"enableTunneling,omitempty"`
	IPConfigurations    *[]BastionHostIPConfiguration           `json:"ipConfigurations,omitempty"`
	NetworkAcls         *BastionHostPropertiesFormatNetworkAcls `json:"networkAcls,omitempty"`
	ProvisioningState   *ProvisioningState                      `json:"provisioningState,omitempty"`
	ScaleUnits          *int64                                  `json:"scaleUnits,omitempty"`
	VirtualNetwork      *SubResource                            `json:"virtualNetwork,omitempty"`
}
