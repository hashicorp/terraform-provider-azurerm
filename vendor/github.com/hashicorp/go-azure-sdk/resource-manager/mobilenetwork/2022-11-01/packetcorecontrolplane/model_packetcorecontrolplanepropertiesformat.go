package packetcorecontrolplane

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCoreControlPlanePropertiesFormat struct {
	ControlPlaneAccessInterface InterfaceProperties                 `json:"controlPlaneAccessInterface"`
	CoreNetworkTechnology       *CoreNetworkType                    `json:"coreNetworkTechnology,omitempty"`
	Installation                *Installation                       `json:"installation,omitempty"`
	InteropSettings             *interface{}                        `json:"interopSettings,omitempty"`
	LocalDiagnosticsAccess      LocalDiagnosticsAccessConfiguration `json:"localDiagnosticsAccess"`
	Platform                    PlatformConfiguration               `json:"platform"`
	ProvisioningState           *ProvisioningState                  `json:"provisioningState,omitempty"`
	RollbackVersion             *string                             `json:"rollbackVersion,omitempty"`
	Sites                       []SiteResourceId                    `json:"sites"`
	Sku                         BillingSku                          `json:"sku"`
	UeMtu                       *int64                              `json:"ueMtu,omitempty"`
	Version                     *string                             `json:"version,omitempty"`
}
