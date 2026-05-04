package loadbalancers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancerVipSwapRequest struct {
	FrontendIPConfigurations *[]LoadBalancerVipSwapRequestFrontendIPConfiguration `json:"frontendIPConfigurations,omitempty"`
}
