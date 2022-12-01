package packetcorecontrolplane

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PacketCoreControlPlanePropertiesFormat struct {
	ControlPlaneAccessInterface InterfaceProperties                  `json:"controlPlaneAccessInterface"`
	CoreNetworkTechnology       *CoreNetworkType                     `json:"coreNetworkTechnology,omitempty"`
	InteropSettings             *interface{}                         `json:"interopSettings,omitempty"`
	LocalDiagnosticsAccess      *LocalDiagnosticsAccessConfiguration `json:"localDiagnosticsAccess,omitempty"`
	MobileNetwork               MobileNetworkResourceId              `json:"mobileNetwork"`
	Platform                    *PlatformConfiguration               `json:"platform,omitempty"`
	ProvisioningState           *ProvisioningState                   `json:"provisioningState,omitempty"`
	Sku                         BillingSku                           `json:"sku"`
	Version                     *string                              `json:"version,omitempty"`
}
