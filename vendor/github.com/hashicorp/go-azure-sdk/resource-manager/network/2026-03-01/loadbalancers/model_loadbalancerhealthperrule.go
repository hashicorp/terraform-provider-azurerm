package loadbalancers

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancerHealthPerRule struct {
	Down                         *int64                                        `json:"down,omitempty"`
	LoadBalancerBackendAddresses *[]LoadBalancerHealthPerRulePerBackendAddress `json:"loadBalancerBackendAddresses,omitempty"`
	Up                           *int64                                        `json:"up,omitempty"`
}
