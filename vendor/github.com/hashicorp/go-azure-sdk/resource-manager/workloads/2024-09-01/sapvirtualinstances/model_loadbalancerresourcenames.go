package sapvirtualinstances

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LoadBalancerResourceNames struct {
	BackendPoolNames             *[]string `json:"backendPoolNames,omitempty"`
	FrontendIPConfigurationNames *[]string `json:"frontendIpConfigurationNames,omitempty"`
	HealthProbeNames             *[]string `json:"healthProbeNames,omitempty"`
	LoadBalancerName             *string   `json:"loadBalancerName,omitempty"`
}
