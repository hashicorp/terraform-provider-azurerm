package virtualwans

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VpnGatewayProperties struct {
	BgpSettings                     *BgpSettings                 `json:"bgpSettings,omitempty"`
	Connections                     *[]VpnConnection             `json:"connections,omitempty"`
	EnableBgpRouteTranslationForNat *bool                        `json:"enableBgpRouteTranslationForNat,omitempty"`
	IPConfigurations                *[]VpnGatewayIPConfiguration `json:"ipConfigurations,omitempty"`
	IsRoutingPreferenceInternet     *bool                        `json:"isRoutingPreferenceInternet,omitempty"`
	NatRules                        *[]VpnGatewayNatRule         `json:"natRules,omitempty"`
	ProvisioningState               *ProvisioningState           `json:"provisioningState,omitempty"`
	VirtualHub                      *SubResource                 `json:"virtualHub,omitempty"`
	VpnGatewayScaleUnit             *int64                       `json:"vpnGatewayScaleUnit,omitempty"`
}
