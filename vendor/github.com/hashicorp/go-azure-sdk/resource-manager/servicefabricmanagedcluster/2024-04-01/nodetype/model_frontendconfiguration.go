package nodetype

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type FrontendConfiguration struct {
	ApplicationGatewayBackendAddressPoolId *string        `json:"applicationGatewayBackendAddressPoolId,omitempty"`
	IPAddressType                          *IPAddressType `json:"ipAddressType,omitempty"`
	LoadBalancerBackendAddressPoolId       *string        `json:"loadBalancerBackendAddressPoolId,omitempty"`
	LoadBalancerInboundNatPoolId           *string        `json:"loadBalancerInboundNatPoolId,omitempty"`
}
