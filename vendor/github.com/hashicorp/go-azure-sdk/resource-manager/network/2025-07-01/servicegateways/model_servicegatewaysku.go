package servicegateways

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceGatewaySku struct {
	Name *ServiceGatewaySkuName `json:"name,omitempty"`
	Tier *ServiceGatewaySkuTier `json:"tier,omitempty"`
}
