package servicegateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceGatewayServicePropertiesFormat struct {
	IsDefault                *bool                 `json:"isDefault,omitempty"`
	LoadBalancerBackendPools *[]BackendAddressPool `json:"loadBalancerBackendPools,omitempty"`
	PublicNatGatewayId       *string               `json:"publicNatGatewayId,omitempty"`
	ServiceType              *ServiceType          `json:"serviceType,omitempty"`
}
