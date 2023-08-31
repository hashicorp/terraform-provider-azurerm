package webapplicationfirewallpolicies

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApplicationGatewayLoadDistributionTargetPropertiesFormat struct {
	BackendAddressPool *SubResource `json:"backendAddressPool,omitempty"`
	WeightPerServer    *int64       `json:"weightPerServer,omitempty"`
}
