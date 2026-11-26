package virtualnetworkgateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VirtualNetworkGatewayNatRuleProperties struct {
	ExternalMappings  *[]VpnNatRuleMapping `json:"externalMappings,omitempty"`
	IPConfigurationId *string              `json:"ipConfigurationId,omitempty"`
	InternalMappings  *[]VpnNatRuleMapping `json:"internalMappings,omitempty"`
	Mode              *VpnNatRuleMode      `json:"mode,omitempty"`
	ProvisioningState *ProvisioningState   `json:"provisioningState,omitempty"`
	Type              *VpnNatRuleType      `json:"type,omitempty"`
}
