package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewaySku struct {
	Capacity *int64                     `json:"capacity,omitempty"`
	Name     *ApplicationGatewaySkuName `json:"name,omitempty"`
	Tier     *ApplicationGatewayTier    `json:"tier,omitempty"`
}
