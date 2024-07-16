package gateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GatewayProperties struct {
	AllowedFeatures   *[]string          `json:"allowedFeatures,omitempty"`
	GatewayEndpoint   *string            `json:"gatewayEndpoint,omitempty"`
	GatewayId         *string            `json:"gatewayId,omitempty"`
	GatewayType       *GatewayType       `json:"gatewayType,omitempty"`
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty"`
}
